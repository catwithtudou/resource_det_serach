package biz

import (
	"context"
	"gorm.io/gorm"
	"resource_det_search/internal/constants"
)

type User struct {
	gorm.Model
	Name     string
	Avatar   string
	Email    string `gorm:"index:idx_email,unique"`
	Pswd     string
	Intro    string
	Role     string
	Sex      string
	School   string
	Sid      string `gorm:"index:idx_sid,unique"`
	IsActive bool
}

type IUserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserBySid(ctx context.Context, sid string) (*User, error)
	GetUserById(ctx context.Context, id uint) (*User, error)
	ListUser(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, id uint, user *User) error
	InsertUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
}

type IUserUsecase interface {
	Login(ctx context.Context, email string, pswd string) (string, *User, constants.ErrCode, error)
}
