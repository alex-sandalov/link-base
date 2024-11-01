package http

import (
	v1 "link-base/internal/http/v1"
	"link-base/internal/service"
	"link-base/pkg/auth"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "link-base/docs"
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

// Init initializes a new Gin HTTP engine and sets up routes for the following endpoints:
//
//   - /swagger/*any: Swagger UI
//   - /ping: Returns "pong" to test the server is up.
func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler()))

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	h.initAPI(router)

	return router
}

// initAPI sets up routes for the API endpoints under /api.
//
// It is a thin wrapper around v1.Handler.Init() that initializes the v1 API
// endpoints and sets them up under the /api group.
func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.service, h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
