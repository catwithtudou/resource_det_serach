package biz

import (
	"context"
	"resource_det_search/internal/constants"
)

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

// TODO:修复 for range 影响的排序部分

type IClassDocumentRepo interface {
	InsertDoc(ctx context.Context, docId uint, cd *ClassDocument) error
	SearchAllQuery(ctx context.Context, queryStr string, sortBy string) ([]*ClassDocument, error)
	SearchQueryByPart(ctx context.Context, queryStr string, partName string, sortBy string) ([]*ClassDocument, error)
	UpdateNums(ctx context.Context, docId uint, likeNum uint, scanNum uint, downloadNum uint) error
	UpdateDimensions(ctx context.Context, docId uint, typeStr constants.DmType, oldDmName string, newDmName string) error
}

// TODO:搜索引擎查询分页处理

type IClassDocumentUsecase interface {
	SearchAllQuery(ctx context.Context, queryStr string, partId uint, sortBy string) ([]*ClassDocument, error)
}
