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

func NewServer(logger *zap.SugaredLogger, userService *service.UserService, dimensionService *service.DimensionService, documentService *service.DocumentService) *Server {
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

	r.POST("/resource/user/upload", documentService.UploadUserDocument)
	r.GET("/resource/user/all", documentService.GetUserAllDocs)
	r.GET("/resource/all", documentService.GetAllDocs)
	r.GET("/resource/user/dimension", documentService.GetUserDimensionDocs)
	r.GET("/resource/user/dimension/all", documentService.GetUserAllDimensionDocs)
	r.GET("/resource/dimension", documentService.GetDimensionDocs)
	r.GET("/resource/dimension/all", documentService.GetAllDimensionDocs)
	r.PUT("/resource/like", documentService.AddLikeDoc)
	r.DELETE("/resource", documentService.DeleteUserDoc)
	r.POST("/resource/user/det", documentService.DetUserDoc)

	return &Server{HttpEngine: r}
}
