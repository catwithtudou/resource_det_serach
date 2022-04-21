package usecase

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mime/multipart"
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

func (u *userUsecase) Register(ctx context.Context, user *biz.User) (constants.ErrCode, error) {
	if user == nil {
		return constants.DefaultErr, errors.New("[UserRegister]the user is nil")
	}

	reUser, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && reUser.ID != 0 {
		return constants.UserEmailExist, errors.New("[UserRegister]the user email is exist")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return constants.DefaultErr, fmt.Errorf("[UserRegister]failed to GetUserByEmail:err=[%+v]", err)
	}

	reUser, err = u.repo.GetUserBySid(ctx, user.Sid)
	if err == nil && reUser.ID != 0 {
		return constants.UserSidExist, fmt.Errorf("[UserRegister]the user sid is exist:err=[%+v]", err)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return constants.DefaultErr, fmt.Errorf("[UserRegister]failed to GetUserBySid:err=[%+v]", err)
	}

	err = u.repo.InsertUser(ctx, user)
	if err != nil {
		return constants.DefaultErr, fmt.Errorf("[UserRegister]failed to insert user:err=[%+v]", err)
	}

	return constants.Success, nil
}

func (u *userUsecase) GetUserInfo(ctx context.Context, id uint) (*biz.User, error) {
	if id <= 0 {
		return nil, errors.New("[GetUserInfo]id is nil")
	}

	user, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[GetUserInfo]failed to GetUserById:err=[%+v]", err)
	}

	return user, nil
}
func (u *userUsecase) UpdateUserInfo(ctx context.Context, user *biz.User) error {
	if user == nil {
		return errors.New("[UpdateUserInfo]user is nil")
	}

	err := u.repo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("[UpdateUserInfo]failed to UpdateUser:err=[%+v]", err)
	}

	return nil
}

func (u *userUsecase) UploadUserAvatar(ctx context.Context, id uint, fileData *multipart.FileHeader) (string, error) {
	if id <= 0 || fileData == nil {
		return "", errors.New("[UploadUserAvatar]user or data is nil")
	}

	fileBytes, err := utils.MultipartFileHeaderToBytes(fileData)
	if err != nil {
		return "", fmt.Errorf("[UploadUserAvatar]failed to MultipartFileHeaderToBytes:err=[%+v]", err)
	}

	key, err := utils.UploadByteData(ctx, fileBytes, utils.GenAvatarKey(id, fileData.Filename))
	if err != nil {
		return "", fmt.Errorf("[UploadUserAvatar]failed to UploadByteData:err=[%+v]", err)
	}

	avatarLink := utils.GenFileLink(key)
	err = u.repo.UpdateUser(ctx, &biz.User{
		Model:  gorm.Model{ID: id},
		Avatar: avatarLink,
	})
	if err != nil {
		return "", fmt.Errorf("[UploadUserAvatar]failed to UpdateUser:err=[%+v]", err)

	}

	return avatarLink, nil
}
