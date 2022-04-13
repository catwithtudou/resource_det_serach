package v1

import (
	"resource_det_search/api"
)

type UserLoginReq struct {
	Email string `form:"email" binding:"required" json:"email"`
	Pswd  string `form:"pswd" binding:"required" json:"pswd"`
}

type UserLoginData struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Role   string `json:"role"`
	School string `json:"school"`
	Sid    string `json:"sid"`
}

type UserLoginResp struct {
	api.RespCommon
	Data *UserLoginData `json:"data,omitempty"`
}

type UserRegisterReq struct {
	Name   string `form:"name" binding:"required" json:"name"`
	Email  string `form:"email" binding:"required" json:"email"`
	Pswd   string `form:"pswd" binding:"required" json:"pswd"`
	Role   string `form:"role" binding:"required" json:"role"`
	Sex    string `form:"sex" binding:"required" json:"sex"`
	Sid    string `form:"sid" binding:"required" json:"sid"`
	School string `form:"school" binding:"required" json:"school"`
}

type UserInfoData struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Email  string `json:"email"`
	Intro  string `json:"intro,omitempty"`
	Role   string `json:"role"`
	Sex    string `json:"sex"`
	School string `json:"school"`
	Sid    string `json:"sid"`
}

type UserGetUserInfoResp struct {
	api.RespCommon
	Data *UserInfoData `json:"data,omitempty"`
}

type UpdateUserInfoReq struct {
	Intro string `form:"intro" binding:"required" json:"intro"`
}

type UserAvatarData struct {
	Avatar string `json:"avatar"`
}

type UploadUserAvatarResp struct {
	api.RespCommon
	Data *UserAvatarData `json:"data,omitempty"`
}
