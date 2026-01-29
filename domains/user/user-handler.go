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
	rg.POST("/users", h.Register())
	rg.POST("/users/login", h.Login())
}

func (h *UserHandler) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto RegisterDto
		if err := ctx.ShouldBindJSON(&dto); err != nil {
			responses.BadRequest(ctx, err)
			return
		}

		token, err := h.service.Register(dto.ToUser())
		if err != nil {
			responses.InternalServerError(ctx, err)
			return
		}

		responses.Created(ctx, token)
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto LoginDto
		if err := ctx.ShouldBindJSON(&dto); err != nil {
			responses.BadRequest(ctx, err)
			return
		}

		token, err := h.service.Login(dto.Email, dto.Password)
		if err != nil {
			responses.Unauthorized(ctx, err)
			return
		}

		responses.Ok(ctx, token)
	}
}