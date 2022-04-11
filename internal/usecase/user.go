package usecase

import (
	"context"
	"errors"
	"fmt"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/constants"
	"resource_det_search/internal/utils"
)

type userUsecase struct {
	repo biz.IUserRepo
}

func NewUserUsecase(repo biz.IUserRepo) biz.IUserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) Login(ctx context.Context, email string, pswd string) (string, *biz.User, constants.ErrCode, error) {
	if email == "" || pswd == "" {
		return "", nil, constants.DefaultErr, errors.New("[UserLogin]the email or pswd is null")
	}

	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, constants.UserEmailErr, fmt.Errorf("[UserLogin]repo get the user by Email failed:err=[%+v]", err)
	}

	if user.IsActive == false {
		return "", nil, constants.UserActiveErr, fmt.Errorf("[UserLogin]the user is not active")
	}

	if user.Pswd != pswd {
		return "", nil, constants.UserPswdErr, fmt.Errorf("[UserLogin]the use pswd is wrong")
	}

	jwtToken, err := utils.GenJwtToken(user.ID)
	if err != nil {
		return "", nil, constants.DefaultErr, fmt.Errorf("[UserLogin]failed to gen jwt token:err=[%+v]", err)
	}

	return jwtToken, user, constants.Success, nil
}
