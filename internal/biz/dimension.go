package biz

import (
	"context"
	"gorm.io/gorm"
)

type Dimension struct {
	gorm.Model
	Uid  uint   `gorm:"index:idx_uid_type"`
	Type string `gorm:"index:idx_uid_type"`
	Name string
}

type IDimensionRepo interface {
	GetDimensionById(ctx context.Context, id uint) (*Dimension, error)
}

type IDimensionUsecase interface {
	GetUserDimension(ctx context.Context, uid uint) ([]*Dimension, error)
}
