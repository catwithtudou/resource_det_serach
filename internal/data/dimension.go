package data

import (
	"context"
	"resource_det_search/internal/biz"
)

type dimensionRepo struct {
	data *Data
}

func NewDimensionRepo(data *Data) biz.IDimensionRepo {
	return &dimensionRepo{
		data: data,
	}
}

func (d *dimensionRepo) GetDimensionById(ctx context.Context, id uint) (*biz.Dimension, error) {
	return nil, nil
}
