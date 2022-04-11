package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"resource_det_search/internal/biz"
)

type DimensionService struct {
	log  *zap.SugaredLogger
	user biz.IDimensionUsecase
}

func NewDimensionService(dimension biz.IDimensionUsecase, logger *zap.SugaredLogger) *DimensionService {
	return &DimensionService{
		log:  logger,
		user: dimension,
	}
}

func (d *DimensionService) GetUserDimension(c *gin.Context) {

}
