package v1

import (
	"link-base/internal/service"
	"link-base/pkg/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service      *service.Service
	tokenManager auth.TokenManager
}

func NewHandler(service *service.Service, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		service:      service,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRouter(v1)
	}
}
