package v1

import "resource_det_search/api"

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
