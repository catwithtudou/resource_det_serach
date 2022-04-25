package biz

import (
	"context"
	"gorm.io/gorm"
	"mime/multipart"
	"resource_det_search/internal/constants"
)

type Document struct {
	gorm.Model
	Uid          uint   `gorm:"not null;index:idx_uid"`
	Type         string `gorm:"not null;size:50"`
	Dir          string `gorm:"not null;size:256"`
	Name         string `gorm:"not null;size:100"`
	Intro        string `gorm:"default:'';size:256"`
	Title        string `gorm:"not null;size:100"`
	DownloadNum  uint   `gorm:"not null;default:0"`
	ScanNum      uint   `gorm:"not null;default:0"`
	LikeNum      uint   `gorm:"not null;default:0"`
	IsLoadSearch bool   `gorm:"not null;default:false"`
	IsSave       bool   `gorm:"not null;default:false"`
	Content      string `gorm:"not null;type:text"`
}

type DocWithDm struct {
	gorm.Model
	DocId uint `gorm:"not null;index:idx_docId"`
	Did   uint `gorm:"not null;index:idx_did"`
}

type IDocumentRepo interface {
	InsertDocWithDms(ctx context.Context, doc *Document, dmIds []uint) (uint, error)
	GetDocsWithDms(ctx context.Context, offset uint, size uint, sortBy string) ([]*Document, map[uint][]*Dimension, error)
	GetDocById(ctx context.Context, id uint) (*Document, error)
	GetDocsByUid(ctx context.Context, uid uint, offset uint, size uint) ([]*Document, map[uint][]*Dimension, error)
	GetDocsWithDid(ctx context.Context, did uint, offset uint, size uint) ([]*Document, error)
	UpdateDocById(ctx context.Context, doc *Document) error
	DeleteDocWithDmsById(ctx context.Context, id uint) error
	DeleteDocWithDmsByIdWithUid(ctx context.Context, id, uid uint) error
	GetSaveDocWithNameAndTitle(ctx context.Context, uid uint, title string) error
	GetDocWithDms(ctx context.Context, id uint) (*Document, []*Dimension, error)
	GetDocsByDidWithDms(ctx context.Context, did uint, offset uint, size uint, sortBy string) ([]*Document, map[uint][]*Dimension, error)
	AddDocNum(ctx context.Context, id uint, num uint, typeStr string) error
	GetUserDocCount(ctx context.Context, uid uint) (int64, error)
	GetDocIdsByDid(ctx context.Context, did uint) ([]uint, error)
}

type IDocumentUsecase interface {
	UploadUserDocument(ctx context.Context, doc *Document, part uint, categories []uint, tags []uint, fileData *multipart.FileHeader) (constants.ErrCode, error)
	GetUserAllDocs(ctx context.Context, uid uint, offset uint, size uint) ([]*Document, map[uint]map[string][]*Dimension, error)
	GetAllDocs(ctx context.Context, offset uint, size uint, sortBy string) ([]*Document, map[uint]map[string][]*Dimension, error)
	GetDmDocs(ctx context.Context, uid uint, did uint, offset uint, size uint) ([]*Document, *Dimension, error)
	GetAllDmTypeDocs(ctx context.Context, uid uint, typeStr string, offset uint, size uint) (map[string][]*Document, error)
	AddLikeDoc(ctx context.Context, docId, num uint) error
	DeleteUserDoc(ctx context.Context, docId uint, uid uint) error
	DetFile(ctx context.Context, fileType string, fileData *multipart.FileHeader) (string, error)
	GetDocWithDms(ctx context.Context, docId uint) (*Document, map[string][]*Dimension, error)
	GetPartDocs(ctx context.Context, did uint, offset uint, size uint, sortBy string) ([]*Document, map[uint]map[string][]*Dimension, error)
	GetUserDocCount(ctx context.Context, uid uint) (int64, error)
}

// TODO:优化：点赞用户限制
// TODO:修复 for range 影响的排序部分
