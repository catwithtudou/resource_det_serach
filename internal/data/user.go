package data

import (
	"context"
	"errors"
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
	if email == "" {
		return nil, errors.New("email is nil")
	}
	result := &biz.User{}
	if err := u.data.db.Model(&biz.User{}).Where("email = ?", email).First(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (u *userRepo) GetUserBySid(ctx context.Context, sid string) (*biz.User, error) {
	if sid == "" {
		return nil, errors.New("sid is nil")
	}
	result := &biz.User{}
	if err := u.data.db.Model(&biz.User{}).Where("sid = ?", sid).First(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (u *userRepo) GetUserById(ctx context.Context, id uint) (*biz.User, error) {
	if id <= 0 {
		return nil, errors.New("id is illegal")
	}
	result := &biz.User{}
	if err := u.data.db.Model(&biz.User{}).Where("id = ?", id).First(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (u *userRepo) ListUser(ctx context.Context) ([]*biz.User, error) {
	result := make([]*biz.User, 0)
	if err := u.data.db.Model(&biz.User{}).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (u *userRepo) UpdateUser(ctx context.Context, user *biz.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if err := u.data.db.Model(&user).Updates(biz.User{Avatar: user.Avatar, Intro: user.Intro, IsActive: user.IsActive}).Error; err != nil {
		return err
	}

	return nil

}
func (u *userRepo) InsertUser(ctx context.Context, user *biz.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if err := u.data.db.Create(user).Error; err != nil {
		return err
	}

	return nil

}
func (u *userRepo) DeleteUser(ctx context.Context, id uint) error {
	if id <= 0 {
		return errors.New("id is illegal")
	}

	if err := u.data.db.Delete(&biz.User{}, id).Error; err != nil {
		return err
	}

	return nil
}
