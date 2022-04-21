package biz

import (
	"context"
	"gorm.io/gorm"
)

type Dimension struct {
	gorm.Model
	Uid  uint   `gorm:"not null;index:idx_uid_type,idx_uid"`
	Type string `gorm:"not null;size:50;index:idx_uid_type"`
	Name string `gorm:"not null;size:50"`
}

type IDimensionRepo interface {
	GetDmById(ctx context.Context, did uint) (*Dimension, error)
	GetDmsByUid(ctx context.Context, uid uint) ([]*Dimension, error)
	GetDmByDidUid(ctx context.Context, did, uid uint) (*Dimension, error)
	InsertDm(ctx context.Context, dm *Dimension) error
	UpdateDm(ctx context.Context, dm *Dimension) error
	DeleteDm(ctx context.Context, did uint) error
	GetDmByUidTypeName(ctx context.Context, uid uint, typeStr string, name string) (*Dimension, error)
	GetDmsByType(ctx context.Context, uid uint, typeStr string) ([]*Dimension, error)
	GetUidsInIds(ctx context.Context, ids []uint) ([]uint, error)
	GetUidTypeInIds(ctx context.Context, ids []uint) ([]*Dimension, error)
	GetDmsInIds(ctx context.Context, ids []uint) ([]*Dimension, error)
	GetDmsPartType(ctx context.Context) ([]*Dimension, error)
}

type IDimensionUsecase interface {
	GetUserDm(ctx context.Context, uid uint) (map[string][]*Dimension, error)
	AddUserDm(ctx context.Context, dm *Dimension) error
	UpdateUserDm(ctx context.Context, did uint, name string, uid uint) error
	DeleteUserDm(ctx context.Context, did, uid uint) error
	GetDmsPartType(ctx context.Context) ([]*Dimension, error)
}
