package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"resource_det_search/internal/service"
)

var ProvideSet = wire.NewSet(NewServer)

type Server struct {
	HttpEngine *gin.Engine
}

func NewServer(userService *service.UserService, dimensionService *service.DimensionService) *Server {
	r := gin.Default()
	r.POST("/user/login", userService.Login)
	r.POST("/user", userService.Register)
	return &Server{HttpEngine: r}
}
