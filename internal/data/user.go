package data

import (
	"context"
	"resource_det_search/internal/biz"
)

type userRepo struct {
	data *Data
}

func NewUserRepo(data *Data) biz.IUserRepo {
	return &userRepo{
		data: data,
	}
}

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	return nil, nil
}
func (u *userRepo) GetUserBySid(ctx context.Context, sid string) (*biz.User, error) {
	return nil, nil

}
func (u *userRepo) GetUserById(ctx context.Context, id uint) (*biz.User, error) {
	return nil, nil

}
func (u *userRepo) ListUser(ctx context.Context) ([]*biz.User, error) {
	return nil, nil

}
func (u *userRepo) UpdateUser(ctx context.Context, id uint, user *biz.User) error {
	return nil

}
func (u *userRepo) InsertUser(ctx context.Context, user *biz.User) error {
	return nil

}
func (u *userRepo) DeleteUser(ctx context.Context, id uint) error {
	return nil
}
