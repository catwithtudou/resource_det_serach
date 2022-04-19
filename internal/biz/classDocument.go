package biz

import "context"

type ClassDocument struct {
	Id          uint     `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Intro       string   `json:"intro"`
	Part        string   `json:"part"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	FileType    string   `json:"file_type"`
	Username    string   `json:"username"`
	UploadDate  int64    `json:"upload_date"`
	DownloadNum int64    `json:"download_num"`
	ScanNum     int64    `json:"scan_num"`
	LikeNum     int64    `json:"like_num"`
}

type IClassDocumentRepo interface {
	InsertDoc(ctx context.Context, docId uint, cd *ClassDocument) error
}
