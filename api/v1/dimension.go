package v1

import "resource_det_search/api"

type DimensionAddUserDmReq struct {
	Type string `form:"type" binding:"required" json:"type"`
	Name string `form:"name" binding:"required" json:"name"`
}

type DimensionGetUserDmReq struct {
	Uid uint `form:"uid" json:"uid"`
}

type DimensionUserDmData struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type DimensionGetUserDmResp struct {
	api.RespCommon
	Data map[string][]*DimensionUserDmData `json:"data,omitempty"`
}

type DimensionUpdateUserDmReq struct {
	Did  uint   `form:"did" binding:"required" json:"did"`
	Name string `form:"name" binding:"required" json:"name"`
}

type DimensionDeleteUserDmReq struct {
	Did uint `form:"did" binding:"required" json:"did"`
}
