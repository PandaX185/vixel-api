package image

import (
	"errors"
	"vixel/shared/middlewares"
	"vixel/shared/responses"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	imageService  *ImageService
	uploadService *UploadService
}

func NewImageHandler(service *ImageService, uploadService *UploadService) *ImageHandler {
	return &ImageHandler{imageService: service, uploadService: uploadService}
}

func (h *ImageHandler) SetupImageRoutes(rg *gin.RouterGroup) {
	rg.POST("/images", middlewares.JWTMiddleware(), h.UploadImage())
	// rg.GET("/images/:id", h.GetImage())
	// rg.GET("/users/:user_id/images", h.ListUserImages())
	// rg.DELETE("/images/:id", h.DeleteImage())
}

func (h *ImageHandler) UploadImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto SaveImageDto
		if err := ctx.ShouldBind(&dto); err != nil {
			responses.BadRequest(ctx, err)
			return
		}

		if errMsg := dto.IsValid(); errMsg != "" {
			responses.BadRequest(ctx, errors.New(errMsg))
			return
		}

		imageURL, err := h.uploadService.UploadImage(dto.File)
		if err != nil {
			responses.InternalServerError(ctx, err)
			return
		}

		image := &Image{
			UserID:  ctx.Value("user_id").(uint),
			URL:     imageURL,
			AltText: dto.AltText,
		}
		savedImage, err := h.imageService.SaveImage(image)
		if err != nil {
			responses.InternalServerError(ctx, err)
			return
		}

		response := ImageResponse{
			ID:      savedImage.ID,
			URL:     savedImage.URL,
			AltText: savedImage.AltText,
			UserID:  savedImage.UserID,
		}

		responses.Created(ctx, response)
	}
}
