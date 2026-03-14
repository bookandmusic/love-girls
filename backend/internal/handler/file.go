package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/internal/service"
)

type FileHandler struct {
	Service *service.FileService
}

func NewFileHandler(service *service.FileService) *FileHandler {
	return &FileHandler{
		Service: service,
	}
}

// RegisterRoutes 注册文件相关的路由
func (h *FileHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	fileGroup := apiGroup.Group("/file")
	{
		fileGroup.POST("/upload", h.SaveFile)
		fileGroup.GET("/:id", h.GetFile)
	}
}

type FileUploadResponse struct {
	FileID uint64                `json:"id"`
	File   *service.FileResponse `json:"file"`
}

type SaveFileRequest struct {
	File            *multipart.FileHeader `form:"file" binding:"required"`
	Path            string                `form:"path"`
	Hash            string                `form:"hash" binding:"required"`
	ThumbnailWidth  int                   `form:"thumbnailWidth"  binding:"omitempty,gte=0,lte=4096"`
	ThumbnailHeight int                   `form:"thumbnailHeight" binding:"omitempty,gte=0,lte=4096"`
}

// SaveFile uploads a file with optional metadata.
// @Summary Upload a file
// @Description Uploads a file to the server. Supports optional fields like filename, path, mimeType, hash, and thumbnail dimensions.
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param path formData string false "Storage path prefix"
// @Param hash formData string false "File content hash (e.g., SHA256)"
// @Param thumbnailWidth formData string false "Desired thumbnail width (for image processing)"
// @Param thumbnailHeight formData string false "Desired thumbnail height (for image processing)"
// @Success 200 {object} Response{data=FileUploadResponse} "File uploaded successfully"
// @Failure 400 {object} Response "Missing or invalid file"
// @Failure 500 {object} Response "Internal server error during file saving or URL generation"
// @Router /file/upload [post]
func (h *FileHandler) SaveFile(c *gin.Context) {
	// 手动解析multipart表单，设置更大的内存限制
	const maxMemory = 200 << 20 // 200 MB
	if err := c.Request.ParseMultipartForm(maxMemory); err != nil {
		h.Service.Log.Error("解析multipart表单失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "文件上传失败",
		})
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.Service.Log.Error("获取上传文件失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "请选择文件",
		})
		return
	}
	defer file.Close()

	// 获取其他表单字段
	hash := c.Request.FormValue("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "文件哈希值不能为空",
		})
		return
	}

	path := c.Request.FormValue("path")
	filename := header.Filename
	size := header.Size
	mimeType := header.Header.Get("Content-Type")

	// 调用 Service
	savedFile, err := h.Service.SaveFile(c, filename, path, mimeType, hash, size, file)
	if err != nil {
		h.Service.Log.Error("文件保存失败", "filename", filename, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
		})
		return
	}

	// 返回
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "文件保存成功",
		Data: FileUploadResponse{
			FileID: savedFile.ID,
			File:   h.Service.BuildFileResponse(c, savedFile),
		},
	})
}

// GetFile retrieves a file by its ID and streams it back to the client.
// @Summary Get a file by ID
// @Description Returns the file content as a stream. The browser will either preview or download based on Content-Disposition.
// @Tags files
// @Produce application/octet-stream
// @Param id path string true "File ID"
// @Success 200 {file} file "File stream"
// @Failure 400 {object} Response "Invalid file ID format"
// @Failure 500 {object} Response "Failed to read file from storage"
// @Router /file/{id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	idStr := c.Param("id")

	var id uint64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		h.Service.Log.Error("无效的文件ID格式", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的文件ID格式",
		})
		return
	}

	width, _ := strconv.Atoi(c.Query("w"))

	// 需要缩略图且有 ImageProxy 内网地址 -> 代理到 ImageProxy
	if width > 0 && h.Service.HasImageProxyInternal() {
		h.proxyToImageProxy(c, id, width)
		return
	}

	// 直接返回文件
	fileReader, file, err := h.Service.ReadFile(c, id)
	if err != nil {
		h.Service.Log.Error("文件读取失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
		})
		return
	}
	defer fileReader.Close()

	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, url.QueryEscape(file.OriginalName)))

	_, err = io.Copy(c.Writer, fileReader)
	if err != nil {
		h.Service.Log.Error("文件流复制失败", "id", id, "error", err)
		return
	}
}

// proxyToImageProxy 代理请求到 ImageProxy
func (h *FileHandler) proxyToImageProxy(c *gin.Context, fileID uint64, width int) {
	proxyURL := h.Service.GetImageProxyURL(c, fileID, width)
	if proxyURL == "" {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
		})
		return
	}

	resp, err := http.Get(proxyURL)
	if err != nil {
		h.Service.Log.Error("ImageProxy 请求失败", "url", proxyURL, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
		})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Content-Disposition", resp.Header.Get("Content-Disposition"))

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		h.Service.Log.Error("ImageProxy 响应流复制失败", "error", err)
		return
	}
}
