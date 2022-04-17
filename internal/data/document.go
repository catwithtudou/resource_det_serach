package data

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"resource_det_search/internal/biz"
)

type documentRepo struct {
	data *Data
}

func NewDocumentRepo(data *Data) biz.IDocumentRepo {
	return &documentRepo{
		data: data,
	}
}

func (d *documentRepo) GetDocs(ctx context.Context) ([]*biz.Document, error) {
	docs := make([]*biz.Document, 0)
	if err := d.data.db.Model(&biz.Document{}).Find(&docs).Error; err != nil {
		return nil, err
	}

	return docs, nil
}
func (d *documentRepo) GetDocById(ctx context.Context, id uint) (*biz.Document, error) {
	if id <= 0 {
		return nil, errors.New("id is nil")
	}
	result := &biz.Document{}
	if err := d.data.db.Model(&biz.Document{}).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (d *documentRepo) GetDocsByUid(ctx context.Context, uid uint) ([]*biz.Document, error) {
	if uid <= 0 {
		return nil, errors.New("uid is nil")
	}

	result := make([]*biz.Document, 0)
	if err := d.data.db.Model(&biz.Document{}).Where("uid = ?", uid).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (d *documentRepo) GetDocsWithDid(ctx context.Context, did uint) ([]*biz.Document, error) {
	if did <= 0 {
		return nil, errors.New("did is nil")
	}

	result := make([]*biz.Document, 0)
	subQuery := d.data.db.Select("doc_id").Where("did = ?", did).Table("doc_with_dm")
	if err := d.data.db.Model(&biz.Document{}).Where("id in (?)", subQuery).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (d *documentRepo) UpdateDocById(ctx context.Context, doc *biz.Document) error {
	if doc == nil {
		return errors.New("doc is nil")
	}

	return d.data.db.Model(&doc).Updates(biz.Document{
		Intro: doc.Intro,
		Title: doc.Title,
		//DownloadNum:  doc.DownloadNum,
		//ScanNum:      doc.ScanNum,
		//LikeNum:      doc.LikeNum,
		IsLoadSearch: doc.IsLoadSearch,
		IsSave:       doc.IsSave,
		Content:      doc.Content,
	}).Error
}
func (d *documentRepo) AddDocLikeNum(ctx context.Context, id uint, num uint) error {
	if id <= 0 || num <= 0 {
		return errors.New("id or num is nil")
	}

	if err := d.data.db.Model(&biz.Document{}).Where("id = ?", id).First(&biz.Document{}).Error; err != nil {
		return err
	}

	return d.data.db.Model(&biz.Document{
		Model: gorm.Model{ID: id},
	}).UpdateColumn("like_num", gorm.Expr("like_num + ?", num)).Error
}
func (d *documentRepo) DeleteDocById(ctx context.Context, id uint) error {
	if id <= 0 {
		return errors.New("id is nil")
	}

	//同时删除doc_with_dm
	tx := d.data.db.Begin()

	if err := tx.Delete(&biz.Document{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&biz.DocWithDm{}, "doc_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
func (d *documentRepo) DeleteDocByIdWithUid(ctx context.Context, id, uid uint) error {
	if id <= 0 || uid < 0 {
		return errors.New("id or uid is nil")
	}

	//同时删除doc_with_dm
	tx := d.data.db.Begin()

	if err := tx.Model(&biz.Document{}).Where("id = ? and uid = ?", id, uid).First(&biz.Document{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&biz.Document{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&biz.DocWithDm{}, "doc_id = ? ", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
