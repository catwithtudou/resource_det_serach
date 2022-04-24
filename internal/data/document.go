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

func (d *documentRepo) InsertDocWithDms(ctx context.Context, doc *biz.Document, dmIds []uint) (uint, error) {
	if doc == nil {
		return 0, errors.New("doc is nil")
	}

	if len(dmIds) == 0 {
		return 0, d.data.db.Create(doc).Error
	}

	tx := d.data.db.Begin()

	if err := d.data.db.Create(doc).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	docWithDms := make([]*biz.DocWithDm, 0, len(dmIds))
	for _, v := range dmIds {
		docWithDms = append(docWithDms, &biz.DocWithDm{
			DocId: doc.ID,
			Did:   v,
		})
	}

	if err := d.data.db.Create(&docWithDms).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return doc.ID, nil
}
func (d *documentRepo) GetDocsWithDms(ctx context.Context, offset uint, size uint) ([]*biz.Document, map[uint][]*biz.Dimension, error) {
	if offset < 0 || size <= 0 {
		return nil, nil, errors.New("offset or size is nil")
	}

	docs := make([]*biz.Document, 0)
	if err := d.data.db.Model(&biz.Document{}).Select("id,created_at,updated_at,uid,type,name,intro,title,download_num,scan_num,like_num,is_load_search,is_save").Order("id").Limit(int(size)).Offset(int(offset)).Find(&docs).Error; err != nil {
		return nil, nil, err
	}

	resDms := make(map[uint][]*biz.Dimension)
	for _, v := range docs {
		dms := make([]*biz.Dimension, 0)
		subQuery := d.data.db.Select("did").Where("doc_id = ?", v.ID).Table("doc_with_dm")
		err := d.data.db.Model(&biz.Dimension{}).Where("id in (?)", subQuery).Find(&dms).Error
		if err != nil {
			return nil, nil, err
		}
		resDms[v.ID] = dms
	}

	return docs, resDms, nil
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
func (d *documentRepo) GetDocsByUid(ctx context.Context, uid uint, offset uint, size uint) ([]*biz.Document, map[uint][]*biz.Dimension, error) {
	if uid <= 0 || offset < 0 || size <= 0 {
		return nil, nil, errors.New("uid is nil")
	}

	result := make([]*biz.Document, 0)
	if err := d.data.db.Model(&biz.Document{}).Select("id,created_at,updated_at,uid,type,name,intro,title,download_num,scan_num,like_num,is_load_search,is_save").Where("uid = ?", uid).Order("id").Limit(int(size)).Offset(int(offset)).Find(&result).Error; err != nil {
		return nil, nil, err
	}

	resDms := make(map[uint][]*biz.Dimension)
	for _, v := range result {
		dms := make([]*biz.Dimension, 0)
		subQuery := d.data.db.Select("did").Where("doc_id = ?", v.ID).Table("doc_with_dm")
		err := d.data.db.Model(&biz.Dimension{}).Where("id in (?)", subQuery).Find(&dms).Error
		if err != nil {
			return nil, nil, err
		}
		resDms[v.ID] = dms
	}

	return result, resDms, nil
}
func (d *documentRepo) GetDocsWithDid(ctx context.Context, did uint, offset uint, size uint) ([]*biz.Document, error) {
	if did <= 0 || offset < 0 || size <= 0 {
		return nil, errors.New("did or offset or size is nil")
	}

	result := make([]*biz.Document, 0)
	subQuery := d.data.db.Select("doc_id").Where("did = ?", did).Table("doc_with_dm")
	if err := d.data.db.Model(&biz.Document{}).Select("id,created_at,updated_at,uid,type,name,intro,title,download_num,scan_num,like_num,is_load_search,is_save").Where("id in (?)", subQuery).Order("id").Limit(int(size)).Offset(int(offset)).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (d *documentRepo) UpdateDocById(ctx context.Context, doc *biz.Document) error {
	if doc == nil {
		return errors.New("doc is nil")
	}

	return d.data.db.Model(&doc).Updates(biz.Document{
		Dir:   doc.Dir,
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
func (d *documentRepo) DeleteDocWithDmsById(ctx context.Context, id uint) error {
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
func (d *documentRepo) DeleteDocWithDmsByIdWithUid(ctx context.Context, id, uid uint) error {
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
func (d *documentRepo) GetSaveDocWithNameAndTitle(ctx context.Context, uid uint, title string) error {
	if uid <= 0 || title == "" {
		return errors.New("name or title is nil")
	}

	return d.data.db.Model(&biz.Document{}).Where("uid = ? AND is_save = 1 AND title = ?", uid, title).First(&biz.Document{}).Error
}

func (d *documentRepo) GetDocWithDms(ctx context.Context, id uint) (*biz.Document, []*biz.Dimension, error) {
	if id <= 0 {
		return nil, nil, errors.New("id is nil")
	}

	doc, err := d.GetDocById(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	dms := make([]*biz.Dimension, 0)
	subQuery := d.data.db.Select("did").Where("doc_id = ?", id).Table("doc_with_dm")
	err = d.data.db.Model(&biz.Dimension{}).Where("id in (?)", subQuery).Find(&dms).Error
	if err != nil {
		return nil, nil, err
	}

	return doc, dms, nil
}

func (d *documentRepo) GetDocsByDidWithDms(ctx context.Context, did uint, offset uint, size uint) ([]*biz.Document, map[uint][]*biz.Dimension, error) {
	if did <= 0 || offset < 0 || size <= 0 {
		return nil, nil, errors.New("did or offset or size is nil")
	}

	docs := make([]*biz.Document, 0)
	subQuery := d.data.db.Select("doc_id").Where("did = ?", did).Table("doc_with_dm")
	if err := d.data.db.Model(&biz.Document{}).Select("id,created_at,updated_at,uid,type,name,intro,title,download_num,scan_num,like_num,is_load_search,is_save").Where("id in (?)", subQuery).Order("id").Limit(1).Limit(int(size)).Offset(int(offset)).Find(&docs).Error; err != nil {
		return nil, nil, err
	}

	resDms := make(map[uint][]*biz.Dimension)
	for _, v := range docs {
		dms := make([]*biz.Dimension, 0)
		subQuery := d.data.db.Select("did").Where("doc_id = ?", v.ID).Table("doc_with_dm")
		err := d.data.db.Model(&biz.Dimension{}).Where("id in (?)", subQuery).Find(&dms).Error
		if err != nil {
			return nil, nil, err
		}
		resDms[v.ID] = dms
	}

	return docs, resDms, nil

}

func (d *documentRepo) AddDocNum(ctx context.Context, id uint, num uint, typeStr string) error {
	if id <= 0 || num <= 0 || typeStr == "" {
		return errors.New("id or num is nil")
	}

	if err := d.data.db.Model(&biz.Document{}).Where("id = ?", id).First(&biz.Document{}).Error; err != nil {
		return err
	}

	return d.data.db.Model(&biz.Document{
		Model: gorm.Model{ID: id},
	}).UpdateColumn(typeStr, gorm.Expr(typeStr+" + ?", num)).Error
}
