package biz

import (
	"context"
	"gorm.io/gorm"
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
	GetDocs(ctx context.Context) ([]*Document, error)
	GetDocById(ctx context.Context, id uint) (*Document, error)
	GetDocsByUid(ctx context.Context, uid uint) ([]*Document, error)
	GetDocsWithDid(ctx context.Context, did uint) ([]*Document, error)
	UpdateDocById(ctx context.Context, doc *Document) error
	AddDocLikeNum(ctx context.Context, id uint, num uint) error
	DeleteDocById(ctx context.Context, id uint) error
	DeleteDocByIdWithUid(ctx context.Context, id, uid uint) error
}

type IDocumentUsecase interface {
	GetUserAllDocs(ctx context.Context, uid uint) ([]*Document, error)
	GetAllDocs(ctx context.Context) ([]*Document, error)
	GetDmDocs(ctx context.Context, uid uint, did uint) ([]*Document, *Dimension, error)
	GetAllDmTypeDocs(ctx context.Context, uid uint, typeStr string) (map[string][]*Document, error)
	AddLikeDoc(ctx context.Context, docId, num uint) error
	DeleteUserDoc(ctx context.Context, docId uint, uid uint) error
}
