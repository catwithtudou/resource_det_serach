package biz

import (
	"context"
	"gorm.io/gorm"
	"resource_det_search/internal/constants"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null;size:50"`
	Avatar   string `gorm:"default:'';size:256"`
	Email    string `gorm:"not null;index:idx_email,unique;size:50"`
	Pswd     string `gorm:"not null;size:50"`
	Intro    string `gorm:"default:'';size:256"`
	Role     string `gorm:"not null;size:50"`
	Sex      string `gorm:"not null;size:50"`
	School   string `gorm:"not null;size:50"`
	Sid      string `gorm:"not null;index:idx_sid,unique;size:50"`
	IsActive bool   `gorm:"not null;default:false"`
}

type IUserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserBySid(ctx context.Context, sid string) (*User, error)
	GetUserById(ctx context.Context, id uint) (*User, error)
	ListUser(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	InsertUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
}

type IUserUsecase interface {
	Login(ctx context.Context, email string, pswd string) (string, *User, constants.ErrCode, error)
	Register(ctx context.Context, user *User) (constants.ErrCode, error)
}
