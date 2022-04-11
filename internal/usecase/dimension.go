package usecase

import (
	"context"
	"resource_det_search/internal/biz"
)

type dimensionUsecase struct {
	repo biz.IDimensionRepo
}

func NewDimensionUsecase(repo biz.IDimensionRepo) biz.IDimensionUsecase {
	return &dimensionUsecase{
		repo: repo,
	}
}

func (d *dimensionUsecase) GetUserDimension(ctx context.Context, uid uint) ([]*biz.Dimension, error) {
	return nil, nil
}
