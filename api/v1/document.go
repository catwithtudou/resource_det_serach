package v1

import "resource_det_search/api"

type GetUserAllDocsReq struct {
	Uid    uint `form:"uid" json:"uid"`
	Offset uint `form:"offset" json:"offset"`
	Size   uint `form:"size" binding:"required" json:"size"`
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

type DocNoDmData struct {
	DocId       uint   `json:"doc_id"`
	Uid         uint   `json:"uid"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Intro       string `json:"intro"`
	Title       string `json:"title"`
	DownloadNum uint   `json:"download_num"`
	ScanNum     uint   `json:"scan_num"`
	LikeNum     uint   `json:"like_num"`
	UploadTime  string `json:"upload_time"`
}

type DocPartData struct {
	DocId       uint                  `json:"doc_id"`
	Uid         uint                  `json:"uid"`
	Type        string                `json:"type"`
	Name        string                `json:"name"`
	Intro       string                `json:"intro"`
	Title       string                `json:"title"`
	DownloadNum uint                  `json:"download_num"`
	ScanNum     uint                  `json:"scan_num"`
	LikeNum     uint                  `json:"like_num"`
	UploadTime  string                `json:"upload_time"`
	Part        DimensionUserDmData   `json:"part"`
	Categories  []DimensionUserDmData `json:"categories"`
	Tags        []DimensionUserDmData `json:"tags"`
}

type GetUserAllDocsResp struct {
	api.RespCommon
	Data []*DocPartData `json:"data,omitempty"`
}

type GetAllDocsReq struct {
	Offset uint   `form:"offset" json:"offset"`
	Size   uint   `form:"size" binding:"required" json:"size"`
	SortBy string `form:"sort_by" json:"sort_by"`
}

type GetAllDocsResp struct {
	api.RespCommon
	Data []*DocPartData `json:"data,omitempty"`
}

type GetUserDimensionDocsReq struct {
	Uid    uint `form:"uid" json:"uid"`
	Did    uint `form:"did" binding:"required" json:"did"`
	Offset uint `form:"offset" json:"offset"`
	Size   uint `form:"size" binding:"required" json:"size"`
}

type DimensionDocsData struct {
	Did    uint           `json:"did"`
	DmName string         `json:"dm_name"`
	DmType string         `json:"dm_type"`
	Docs   []*DocNoDmData `json:"docs,omitempty"`
}

type GetUserDimensionDocsResp struct {
	api.RespCommon
	Data *DimensionDocsData `json:"data,omitempty"`
}

type GetUserAllDimensionDocsReq struct {
	Uid    uint   `form:"uid" json:"uid"`
	Type   string `form:"type" binding:"required" json:"type"`
	Offset uint   `form:"offset" json:"offset"`
	Size   uint   `form:"size" binding:"required" json:"size"`
}

type GetUserAllDimensionDocsResp struct {
	api.RespCommon
	Data map[string][]*DocData `json:"data,omitempty"`
}

type GetDimensionDocsReq struct {
	Did    uint   `form:"did" binding:"required" json:"did"`
	Offset uint   `form:"offset" json:"offset"`
	Size   uint   `form:"size" binding:"required" json:"size"`
	SortBy string `form:"sort_by" json:"sort_by"`
}

type GetDimensionDocsResp struct {
	api.RespCommon
	Data []*DocPartData `json:"data,omitempty"`
}

type GetAllDimensionDocsReq struct {
	Type   string `form:"type" binding:"required" json:"type"`
	Offset uint   `form:"offset" json:"offset"`
	Size   uint   `form:"size" binding:"required" json:"size"`
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

type GetUserDocCountReq struct {
	Uid uint `form:"uid" binding:"required" json:"uid"`
}

type UserDocCountData struct {
	Count int64 `json:"count"`
}

type GetUserDocCountResp struct {
	api.RespCommon
	Data *UserDocCountData `json:"data,omitempty"`
}
