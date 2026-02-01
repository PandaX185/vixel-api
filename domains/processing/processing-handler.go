package processing

import (
	"vixel/shared/responses"

	"github.com/gin-gonic/gin"
)

type ProcessingHandler struct {
	processingService *ProcessingService
}

func NewProcessingHandler(service *ProcessingService) *ProcessingHandler {
	return &ProcessingHandler{processingService: service}
}

func (h *ProcessingHandler) SetupProcessingRoutes(rg *gin.RouterGroup) {
	rg.POST("/images/:id/transform", h.TransformImage())
}

func (h *ProcessingHandler) TransformImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto TransformationDTO
		if err := ctx.ShouldBindJSON(&dto); err != nil {
			ctx.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		imageID := ctx.Param("id")
		newImageURL, err := h.processingService.TransformImage(ctx, imageID, dto)
		if err != nil {
			responses.InternalServerError(ctx, err)
			return
		}

		responses.Ok(ctx, gin.H{"new_image_url": newImageURL})
	}
}
