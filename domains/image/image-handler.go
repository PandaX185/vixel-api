package image

import (
	"errors"
	"strconv"
	"vixel/shared/middlewares"
	"vixel/shared/responses"

	"github.com/gin-gonic/gin"
)

type ImageServiceInterface interface {
	SaveImage(image *Image) (*Image, error)
	GetImageByID(id uint) (*Image, error)
	ListImagesByUser(userID uint) ([]Image, error)
	DeleteImage(id uint) error
}

type ImageHandler struct {
	imageService  ImageServiceInterface
	uploadService UploadServiceInterface
}

func NewImageHandler(service ImageServiceInterface, uploadService UploadServiceInterface) *ImageHandler {
	return &ImageHandler{imageService: service, uploadService: uploadService}
}

func (h *ImageHandler) SetupImageRoutes(rg *gin.RouterGroup) {
	rg.POST("/images", middlewares.JWTMiddleware(), h.UploadImage())
	rg.GET("/images/:id", middlewares.JWTMiddleware(), h.GetImage())
	rg.GET("/users/:user_id/images", middlewares.JWTMiddleware(), h.ListUserImages())
	rg.DELETE("/images/:id", middlewares.JWTMiddleware(), h.DeleteImage())
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

		imageURL, err := h.uploadService.UploadImage(ctx, dto.File)
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

func (h *ImageHandler) GetImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			responses.BadRequest(ctx, errors.New("invalid image id"))
			return
		}

		image, err := h.imageService.GetImageByID(uint(id))
		if err != nil {
			responses.NotFound(ctx, errors.New("image not found"))
			return
		}

		userID := ctx.Value("user_id").(uint)
		if image.UserID != userID {
			responses.Unauthorized(ctx, errors.New("access denied"))
			return
		}

		response := ImageResponse{
			ID:      image.ID,
			URL:     image.URL,
			AltText: image.AltText,
			UserID:  image.UserID,
		}

		responses.Ok(ctx, response)
	}
}

func (h *ImageHandler) ListUserImages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr := ctx.Param("user_id")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			responses.BadRequest(ctx, errors.New("invalid user id"))
			return
		}

		images, err := h.imageService.ListImagesByUser(uint(userID))
		if err != nil {
			responses.InternalServerError(ctx, err)
			return
		}

		var imageResponses []ImageResponse
		for _, img := range images {
			imageResponses = append(imageResponses, ImageResponse{
				ID:      img.ID,
				URL:     img.URL,
				AltText: img.AltText,
				UserID:  img.UserID,
			})
		}

		responses.Ok(ctx, imageResponses)
	}
}

func (h *ImageHandler) DeleteImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			responses.BadRequest(ctx, errors.New("invalid image id"))
			return
		}

		image, err := h.imageService.GetImageByID(uint(id))
		if err != nil {
			responses.NotFound(ctx, errors.New("image not found"))
			return
		}

		userID := ctx.Value("user_id").(uint)
		if image.UserID != userID {
			responses.Unauthorized(ctx, errors.New("access denied"))
			return
		}

		err = h.imageService.DeleteImage(uint(id))
		if err != nil {
			responses.InternalServerError(ctx, err)
			return
		}

		responses.Ok(ctx, gin.H{"message": "image deleted"})
	}
}
