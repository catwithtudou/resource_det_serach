package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
	"resource_det_search/internal/service"
)

var ProvideSet = wire.NewSet(NewServer)

type Server struct {
	HttpEngine *gin.Engine
}

func NewServer(logger *zap.SugaredLogger, userService *service.UserService, dimensionService *service.DimensionService) *Server {
	r := gin.Default()
	r.POST("/user/login", userService.Login)
	r.POST("/user", userService.Register)

	r.GET("/user", AuthJwtMw(logger), userService.GetUserInfo)
	r.PUT("/user", AuthJwtMw(logger), userService.UpdateUserInfo)

	return &Server{HttpEngine: r}
}
