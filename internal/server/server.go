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

	r.MaxMultipartMemory = 8 << 20

	r.POST("/user/login", userService.Login)
	r.POST("/user", userService.Register)

	r.Use(AuthJwtMw(logger))

	r.GET("/user", userService.GetUserInfo)
	r.PUT("/user", userService.UpdateUserInfo)
	r.POST("/user/avatar", userService.UploadUserAvatar)

	r.POST("/dimension", dimensionService.AddUserDm)
	r.GET("/dimension", dimensionService.GetUserDm)
	r.PUT("/dimension", dimensionService.UpdateUserDm)
	r.DELETE("/dimension", dimensionService.DeleteUserDm)

	return &Server{HttpEngine: r}
}
