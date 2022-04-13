package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"resource_det_search/api"
	v1 "resource_det_search/api/v1"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/constants"
	"resource_det_search/internal/utils"
)

type UserService struct {
	log  *zap.SugaredLogger
	user biz.IUserUsecase
}

func NewUserService(user biz.IUserUsecase, logger *zap.SugaredLogger) *UserService {
	return &UserService{
		log:  logger,
		user: user,
	}
}

func (u *UserService) Login(c *gin.Context) {
	var req v1.UserLoginReq
	if err := c.Bind(&req); err != nil {
		u.log.Errorf("[UserService-Login]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, v1.UserLoginResp{
			RespCommon: api.FormEmptyErr,
		})
		return
	}

	if len(req.Email) > 100 || len(req.Pswd) > 100 || !utils.CheckEmail(req.Email) {
		u.log.Errorf("[UserService-Login]illegal params")
		c.JSON(http.StatusOK, v1.UserLoginResp{
			RespCommon: api.FormIllegalErr,
		})
		return
	}

	token, user, code, err := u.user.Login(c, req.Email, req.Pswd)
	if err != nil && code == constants.DefaultErr {
		u.log.Errorf("[UserService-Login]failed to login:err=[%+v],code=[%+v]", err, code)
		c.JSON(http.StatusOK, v1.UserLoginResp{RespCommon: api.DefaultErr})
		return
	}
	if err != nil && code == constants.UserEmailErr {
		u.log.Errorf("[UserService-Login]failed to login:err=[%+v],code=[%+v]", err, code)
		c.JSON(http.StatusOK, v1.UserLoginResp{RespCommon: api.UserEmailNotExist})
		return
	}

	if err != nil && code == constants.UserPswdErr {
		u.log.Errorf("[UserService-Login]failed to login:err=[%+v],code=[%+v]", err, code)
		c.JSON(http.StatusOK, v1.UserLoginResp{RespCommon: api.UserPswdErr})
		return
	}

	if err != nil && code == constants.UserActiveErr {
		u.log.Errorf("[UserService-Login]failed to login:err=[%+v],code=[%+v]", err, code)
		c.JSON(http.StatusOK, v1.UserLoginResp{RespCommon: api.UserNotActive})
		return
	}

	c.JSON(http.StatusOK, v1.UserLoginResp{
		RespCommon: api.Success,
		Data: &v1.UserLoginData{
			Token:  token,
			Name:   user.Name,
			Avatar: user.Avatar,
			Role:   user.Role,
			School: user.School,
			Sid:    user.Sid,
		},
	})
	return
}

func (u *UserService) Register(c *gin.Context) {
	var req v1.UserRegisterReq
	if err := c.Bind(&req); err != nil {
		u.log.Errorf("[UserService-Register]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if len(req.Email) > 100 || len(req.School) > 100 || len(req.Name) > 100 || !utils.CheckEmail(req.Email) || !utils.CheckRole(req.Role) || !utils.CheckSex(req.Sex) || utils.CheckPswd(req.Pswd) {
		u.log.Errorf("[UserService-Register]illegal params")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	code, err := u.user.Register(c, &biz.User{
		Name:     req.Name,
		Email:    req.Email,
		Pswd:     req.Pswd,
		Role:     req.Role,
		Sex:      req.Sex,
		School:   req.School,
		Sid:      req.Sid,
		IsActive: true,
	})
	if err != nil && code == constants.DefaultErr {
		u.log.Errorf("[UserService-Register]failed to register:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	if err != nil && code == constants.UserEmailExist {
		u.log.Errorf("[UserService-Register]the user email is exist:err=[%+v]", err)
		c.JSON(http.StatusOK, api.UserEmailExist)
		return
	}

	if err != nil && code == constants.UserSidExist {
		u.log.Errorf("[UserService-Register]the user sid is exist:err=[%+v]", err)
		c.JSON(http.StatusOK, api.UserSidExist)
		return
	}

	c.JSON(http.StatusOK, api.Success)
	return
}

func (u *UserService) GetUserInfo(c *gin.Context) {
	var req v1.UserGetUserInfoReq
	var uid uint
	if _ = c.Bind(&req); req.Uid > 0 {
		uid = req.Uid
	}

	if uid <= 0 {
		getUid, ok := c.Get("uid")
		if !ok || getUid.(uint) <= 0 {
			u.log.Errorf("[UserService-GetUserInfo]failed to get uid")
			c.JSON(http.StatusOK, v1.UserGetUserInfoResp{
				RespCommon: api.UserAuthErr,
			})
			return
		}
		uid = getUid.(uint)
	}

	user, err := u.user.GetUserInfo(c, uid)
	if err != nil {
		u.log.Errorf("[UserService-GetUserInfo]failed to GetUserInfo:err=[%+v]", err)
		c.JSON(http.StatusOK, v1.UserGetUserInfoResp{
			RespCommon: api.DefaultErr,
		})
		return
	}

	c.JSON(http.StatusOK, v1.UserGetUserInfoResp{
		RespCommon: api.Success,
		Data: &v1.UserInfoData{
			Name:   user.Name,
			Avatar: user.Avatar,
			Email:  user.Email,
			Intro:  user.Intro,
			Role:   user.Role,
			Sex:    user.Sex,
			School: user.School,
			Sid:    user.Sid,
		},
	})
	return
}

func (u *UserService) UpdateUserInfo(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok || uid.(uint) <= 0 {
		u.log.Errorf("[UserService-UpdateUserInfo]failed to get uid")
		c.JSON(http.StatusOK, api.UserAuthErr)
		return
	}

	var req v1.UpdateUserInfoReq
	if err := c.Bind(&req); err != nil {
		u.log.Errorf("[UserService-UpdateUserInfo]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if len(req.Intro) > 100 {
		u.log.Errorf("[UserService-UpdateUserInfo]illegal params")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	err := u.user.UpdateUserInfo(c, &biz.User{
		Model: gorm.Model{ID: uid.(uint)},
		Intro: req.Intro,
	})
	if err != nil {
		u.log.Errorf("[UserService-UpdateUserInfo]failed to UpdateUserInfo:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, api.Success)
	return
}

func (u *UserService) UploadUserAvatar(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok || uid.(uint) <= 0 {
		u.log.Errorf("[UserService-UploadUserAvatar]failed to get uid")
		c.JSON(http.StatusOK, api.UserAuthErr)
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		u.log.Errorf("[UserService-UploadUserAvatar]failed to FormFile:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormFileErr)
		return
	}

	link, err := u.user.UploadUserAvatar(c, uid.(uint), file)
	if err != nil {
		u.log.Errorf("[UserService-UploadUserAvatar]failed to UploadUserAvatar:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, v1.UploadUserAvatarResp{
		RespCommon: api.Success,
		Data:       &v1.UserAvatarData{Avatar: link},
	})
	return
}
