package user

import (
	"vixel/shared/responses"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) SetupUserRoutes(rg *gin.RouterGroup) {
	rg.POST("/users/register", h.Register())
}

func (h *UserHandler) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto RegisterDto
		if err := ctx.ShouldBindJSON(&dto); err != nil {
			responses.BadRequest(ctx, err)
			return
		}

		createdUser, err := h.service.Register(dto.ToUser())
		if err != nil {
			responses.InternalServerError(ctx, err)
			return
		}

		responses.Created(ctx, createdUser.ToResponse())
	}
}
