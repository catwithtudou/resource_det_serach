package v1

import "resource_det_search/api"

type GetUserAllDocsReq struct {
	Uid uint `form:"uid" json:"uid"`
}

type DocData struct {
	DocId       uint   `json:"doc_id"`
	Uid         uint   `json:"uid"`
	Type        string `json:"type"`
	Dir         string `json:"dir"`
	Name        string `json:"name"`
	Intro       string `json:"intro"`
	Title       string `json:"title"`
	DownloadNum uint   `json:"download_num"`
	ScanNum     uint   `json:"scan_num"`
	LikeNum     uint   `json:"like_num"`
	Content     string `json:"content"`
}

type GetUserAllDocsResp struct {
	api.RespCommon
	Data []*DocData `json:"data,omitempty"`
}

type GetAllDocsResp struct {
	api.RespCommon
	Data []*DocData `json:"data,omitempty"`
}

type GetUserDimensionDocsReq struct {
	Uid uint `form:"uid" json:"uid"`
	Did uint `form:"did" binding:"required" json:"did"`
}

type DimensionDocsData struct {
	Did    uint       `json:"did"`
	DmName string     `json:"dm_name"`
	DmType string     `json:"dm_type"`
	Docs   []*DocData `json:"docs,omitempty"`
}

type GetUserDimensionDocsResp struct {
	api.RespCommon
	Data *DimensionDocsData `json:"data,omitempty"`
}

type GetUserAllDimensionDocsReq struct {
	Uid  uint   `form:"uid" json:"uid"`
	Type string `form:"type" binding:"required" json:"type"`
}

type GetUserAllDimensionDocsResp struct {
	api.RespCommon
	Data map[string][]*DocData `json:"data,omitempty"`
}

type GetDimensionDocsReq struct {
	Did uint `form:"did" binding:"required" json:"did"`
}

type GetDimensionDocsResp struct {
	api.RespCommon
	Data *DimensionDocsData `json:"data,omitempty"`
}

type GetAllDimensionDocsReq struct {
	Type string `form:"type" binding:"required" json:"type"`
}

type GetAllDimensionDocsResp struct {
	api.RespCommon
	Data map[string][]*DocData `json:"data,omitempty"`
}

type AddLikeDocReq struct {
	DocId uint `form:"doc_id" binding:"required" json:"doc_id"`
}

type DeleteUserDocReq struct {
	DocId uint `form:"doc_id" binding:"required" json:"doc_id"`
}

type UploadUserDocumentReq struct {
	Title      string `form:"title" binding:"required" json:"title"`
	Intro      string `form:"intro" binding:"required" json:"intro"`
	Part       uint   `form:"part" binding:"required" json:"part"`
	Categories string `form:"categories"  json:"categories"`
	Tags       string `form:"tags"  json:"tags"`
}

type DetUserDocResp struct {
	api.RespCommon
	Data string `json:"data"`
}

type GetDocWithDmsReq struct {
	DocId uint `form:"doc_id" json:"doc_id" binding:"required"`
}

type DocWithDmsData struct {
	DocId       uint                  `json:"doc_id"`
	Uid         uint                  `json:"uid"`
	Type        string                `json:"type"`
	Dir         string                `json:"dir"`
	Name        string                `json:"name"`
	Intro       string                `json:"intro"`
	Title       string                `json:"title"`
	DownloadNum uint                  `json:"download_num"`
	ScanNum     uint                  `json:"scan_num"`
	LikeNum     uint                  `json:"like_num"`
	Content     string                `json:"content"`
	Part        DimensionUserDmData   `json:"part"`
	Categories  []DimensionUserDmData `json:"categories"`
	Tags        []DimensionUserDmData `json:"tags"`
}

type GetDocWithDmsResp struct {
	api.RespCommon
	Data *DocWithDmsData `json:"data,omitempty"`
}
