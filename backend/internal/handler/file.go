package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

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
// FileUploadResponse represents the response data for a successful file upload.
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
	path := c.Request.FormValue("path")
	thumbnailWidth := 0
	thumbnailHeight := 0

	// 验证必填字段
	if hash == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "文件哈希值不能为空",
		})
		return
	}

	// ---------- 2. 默认值补齐 ----------
	filename := header.Filename
	size := header.Size

	// ---------- 3. 缩略图参数校验 ----------
	if thumbnailWidth < 0 || thumbnailWidth > 4096 ||
		thumbnailHeight < 0 || thumbnailHeight > 4096 {
		h.Service.Log.Error("缩略图尺寸参数非法", "width", thumbnailWidth, "height", thumbnailHeight)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "缩略图尺寸参数非法",
		})
		return
	}

	// ---------- 5. Optional string 指针 ----------
	mimeType := header.Header.Get("Content-Type")

	// ---------- 6. 调用 Service ----------
	savedFile, err := h.Service.SaveFile(
		c,
		filename,
		path,
		mimeType,
		hash,
		size,
		file,
	)
	if err != nil {
		h.Service.Log.Error("文件保存失败", "filename", filename, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
		})
		return
	}

	// ---------- 7. URL ----------

	// ---------- 8. 返回 ----------
	resp := FileUploadResponse{
		FileID: savedFile.ID,
		File:   h.Service.BuildFileResponse(c, savedFile),
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "文件保存成功",
		Data:    resp,
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

	// 将字符串ID转换为uint类型
	var id uint64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		h.Service.Log.Error("无效的文件ID格式", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的文件ID格式",
			Data:    nil,
		})
		return
	}

	fileReader, file, err := h.Service.ReadFile(c, id)
	if err != nil {
		h.Service.Log.Error("文件读取失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}
	defer fileReader.Close()

	// 手动设置 Header（不设 Content-Length）
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, url.QueryEscape(file.OriginalName)))

	// 直接流式写入，不依赖 Content-Length
	_, err = io.Copy(c.Writer, fileReader)
	if err != nil {
		h.Service.Log.Error("文件流复制失败", "id", id, "error", err)
		return
	}
}
