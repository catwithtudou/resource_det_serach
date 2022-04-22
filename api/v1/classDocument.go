package v1

import "resource_det_search/api"

type SearchAllQueryReq struct {
	Detail string `form:"detail" binding:"required" json:"detail"`
	PartId uint   `form:"part_id" json:"part_id"`
}

type ClassDocumentData struct {
	DocId       uint     `json:"doc_id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Intro       string   `json:"intro"`
	Part        string   `json:"part"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	FileType    string   `json:"file_type"`
	Username    string   `json:"username"`
	UploadTime  string   `json:"upload_time"`
	DownloadNum int64    `json:"download_num"`
	ScanNum     int64    `json:"scan_num"`
	LikeNum     int64    `json:"like_num"`
}

type SearchAllQueryResp struct {
	api.RespCommon
	Data []*ClassDocumentData `json:"data,omitempty"`
}
