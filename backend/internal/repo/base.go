package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FilterCondition struct {
	Field    string      `json:"field"`    // 字段名
	Operator string      `json:"operator"` // 操作符: eq, ne, gt, lt, gte, lte, like, in
	Value    interface{} `json:"value"`    // 值
}

// QueryOptions 查询选项结构体
type QueryOptions struct {
	Conditions   []FilterCondition
	OrderBy      string
	Desc         bool
	Preloads     []string
	PreloadConds map[string][]interface{} // 预加载条件，key为关联名，value为where条件参数
	ForUpdate    bool                     // 是否加锁
}

// QueryOption 函数类型，用于设置查询选项
type QueryOption func(*QueryOptions)

// WithOrder 设置排序字段和方向
func WithOrder(orderBy string, desc bool) QueryOption {
	return func(opts *QueryOptions) {
		opts.OrderBy = "`" + orderBy + "`"
		opts.Desc = desc
	}
}

// WithPreload 设置单个预加载关联
func WithPreload(preload string) QueryOption {
	return func(opts *QueryOptions) {
		opts.Preloads = append(opts.Preloads, preload)
	}
}

// WithPreloads 设置多个预加载关联
func WithPreloads(preloads ...string) QueryOption {
	return func(opts *QueryOptions) {
		opts.Preloads = append(opts.Preloads, preloads...)
	}
}

// WithPreloadCond 设置预加载的过滤条件
func WithPreloadCond(preload string, conds ...interface{}) QueryOption {
	return func(opts *QueryOptions) {
		if opts.PreloadConds == nil {
			opts.PreloadConds = make(map[string][]interface{})
		}
		opts.PreloadConds[preload] = conds
	}
}

// WithForUpdate 设置加锁查询
func WithForUpdate() QueryOption {
	return func(opts *QueryOptions) {
		opts.ForUpdate = true
	}
}

// WithConditions 设置过滤条件
func WithConditions(conditions ...FilterCondition) QueryOption {
	return func(opts *QueryOptions) {
		opts.Conditions = append(opts.Conditions, conditions...)
	}
}

type BaseRepo[T any] struct {
	db *gorm.DB
}

// DB returns the underlying gorm.DB instance for custom queries
func (r *BaseRepo[T]) DB() *gorm.DB {
	return r.db
}

// ApplyFilters applies filter conditions to the database query
// ApplyFilters 应用过滤条件到数据库查询
func (r *BaseRepo[T]) ApplyFilters(db *gorm.DB, conditions []FilterCondition) *gorm.DB {
	for _, condition := range conditions {
		fieldName := "`" + condition.Field + "`"
		switch condition.Operator {
		case "eq":
			db = db.Where(fieldName+" = ?", condition.Value)
		case "ne":
			db = db.Where(fieldName+" != ?", condition.Value)
		case "gt":
			db = db.Where(fieldName+" > ?", condition.Value)
		case "lt":
			db = db.Where(fieldName+" < ?", condition.Value)
		case "gte":
			db = db.Where(fieldName+" >= ?", condition.Value)
		case "lte":
			db = db.Where(fieldName+" <= ?", condition.Value)
		case "like":
			db = db.Where(fieldName+" LIKE ?", "%"+condition.Value.(string)+"%")
		case "in":
			db = db.Where(fieldName+" IN ?", condition.Value)
		}
	}
	return db
}

// NewBaseRepo 创建一个新的基础仓库实例
func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{db: db}
}

// Create 创建实体
func (r *BaseRepo[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update 更新实体
func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// DeleteByID 根据ID删除实体
func (r *BaseRepo[T]) DeleteByID(ctx context.Context, id uint64) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// FindByID 根据ID查找实体，支持预加载关联数据
func (r *BaseRepo[T]) FindByID(ctx context.Context, id uint64, opts ...QueryOption) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx)

	// Apply options
	options := &QueryOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Apply preloads
	for _, preload := range options.Preloads {
		if conds, ok := options.PreloadConds[preload]; ok {
			db = db.Preload(preload, conds...)
		} else {
			db = db.Preload(preload)
		}
	}

	// Apply for update
	if options.ForUpdate {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 查询实体列表，支持过滤、排序、预加载关联数据
func (r *BaseRepo[T]) List(ctx context.Context, opts ...QueryOption) ([]T, error) {
	options := &QueryOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var list []T
	db := r.db.WithContext(ctx)
	db = r.ApplyFilters(db, options.Conditions)

	// Apply preloads
	for _, preload := range options.Preloads {
		if conds, ok := options.PreloadConds[preload]; ok {
			db = db.Preload(preload, conds...)
		} else {
			db = db.Preload(preload)
		}
	}

	// Apply order
	if options.OrderBy != "" {
		if options.Desc {
			db = db.Order(options.OrderBy + " DESC")
		} else {
			db = db.Order(options.OrderBy + " ASC")
		}
	}

	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// FindOne 查找单个实体，支持预加载关联数据
func (r *BaseRepo[T]) FindOne(ctx context.Context, opts ...QueryOption) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx)

	// Apply options
	options := &QueryOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Apply filters
	db = r.ApplyFilters(db, options.Conditions)

	// Apply preloads
	for _, preload := range options.Preloads {
		if conds, ok := options.PreloadConds[preload]; ok {
			db = db.Preload(preload, conds...)
		} else {
			db = db.Preload(preload)
		}
	}

	// Apply for update
	if options.ForUpdate {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := db.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindWithPagination 分页查询实体，支持过滤、排序、预加载关联数据
func (r *BaseRepo[T]) FindWithPagination(ctx context.Context, page, size int, opts ...QueryOption) ([]T, int64, error) {
	var list []T
	var total int64

	// Apply options first to get conditions
	options := &QueryOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db := r.db.WithContext(ctx)
	db = r.ApplyFilters(db, options.Conditions)

	// Count total
	if err := db.Model(&list).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply preloads
	for _, preload := range options.Preloads {
		if conds, ok := options.PreloadConds[preload]; ok {
			db = db.Preload(preload, conds...)
		} else {
			db = db.Preload(preload)
		}
	}

	// Apply order
	if options.OrderBy != "" {
		if options.Desc {
			db = db.Order(options.OrderBy + " DESC")
		} else {
			db = db.Order(options.OrderBy + " ASC")
		}
	}

	// Pagination
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// Count 统计实体总数
func (r *BaseRepo[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountWithConditions 根据条件统计实体总数
func (r *BaseRepo[T]) CountWithConditions(ctx context.Context, conditions ...FilterCondition) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(new(T))
	db = r.ApplyFilters(db, conditions)
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
