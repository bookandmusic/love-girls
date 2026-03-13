package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/love-girl/internal/repo"
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type QueryParams struct {
	Page    int
	Size    int
	SortBy  string
	Order   string
	Filters []repo.FilterCondition
}

var AllowedSortFields = map[string][]string{
	"moments":       {"created_at", "likes"},
	"places":        {"created_at", "name"},
	"wishes":        {"created_at"},
	"anniversaries": {"date", "created_at"},
	"albums":        {"created_at", "name"},
}

var AllowedFilterFields = map[string]map[string][]string{
	"moments": {
		"is_public": {"eq"},
		"user_id":   {"eq"},
		"likes":     {"eq", "gt", "lt", "gte", "lte"},
	},
	"wishes": {
		"approved": {"eq"},
	},
	"albums": {
		"name": {"like"},
	},
	"places": {
		"name": {"like"},
	},
}

func ParseQueryParams(c *gin.Context, resource string) *QueryParams {
	params := &QueryParams{
		Page:  1,
		Size:  DefaultPageSize,
		Order: "desc",
	}

	parsePageParams(c, params)
	parseSortParams(c, params, resource)
	parseFilterParams(c, params, resource)

	return params
}

func parsePageParams(c *gin.Context, params *QueryParams) {
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	if size := c.Query("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil && s > 0 {
			if s > MaxPageSize {
				params.Size = MaxPageSize
			} else {
				params.Size = s
			}
		}
	}
}

func parseSortParams(c *gin.Context, params *QueryParams, resource string) {
	if allowedFields, ok := AllowedSortFields[resource]; ok {
		if sortBy := c.Query("sort_by"); sortBy != "" {
			if isAllowedField(allowedFields, sortBy) {
				params.SortBy = sortBy
			}
		}
	}
	if order := c.Query("order"); order == "asc" || order == "desc" {
		params.Order = order
	}
}

func parseFilterParams(c *gin.Context, params *QueryParams, resource string) {
	filters := c.QueryArray("filter")
	if len(filters) == 0 {
		return
	}
	for _, f := range filters {
		if cond := parseFilter(f, resource); cond != nil {
			params.Filters = append(params.Filters, *cond)
		}
	}
}

func ParsePagination(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = DefaultPageSize
	}
	if size > MaxPageSize {
		size = MaxPageSize
	}
	return page, size
}

func isAllowedField(allowed []string, field string) bool {
	for _, f := range allowed {
		if f == field {
			return true
		}
	}
	return false
}

func parseFilter(filterStr, resource string) *repo.FilterCondition {
	parts := strings.Split(filterStr, ":")
	if len(parts) != 3 {
		return nil
	}

	field, op, value := parts[0], parts[1], parts[2]

	allowedOps, ok := AllowedFilterFields[resource][field]
	if !ok {
		return nil
	}
	if !isAllowedField(allowedOps, op) {
		return nil
	}

	return &repo.FilterCondition{
		Field:    field,
		Operator: op,
		Value:    parseValue(value),
	}
}

func parseValue(value string) interface{} {
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}
	return value
}
