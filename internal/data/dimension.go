package data

import (
	"context"
	"errors"
	"resource_det_search/internal/biz"
)

type dimensionRepo struct {
	data *Data
}

func NewDimensionRepo(data *Data) biz.IDimensionRepo {
	return &dimensionRepo{
		data: data,
	}
}

func (d *dimensionRepo) GetDmById(ctx context.Context, did uint) (*biz.Dimension, error) {
	if did <= 0 {
		return nil, errors.New("did is nil")
	}

	result := &biz.Dimension{}
	if err := d.data.db.Model(&biz.Dimension{}).Where("id = ?", did).First(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (d *dimensionRepo) GetDmsByUid(ctx context.Context, uid uint) ([]*biz.Dimension, error) {
	if uid < 0 {
		return nil, errors.New("uid is nil")
	}

	result := make([]*biz.Dimension, 0)
	if err := d.data.db.Model(&biz.Dimension{}).Where("uid = ?", uid).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (d *dimensionRepo) GetDmByDidUid(ctx context.Context, did, uid uint) (*biz.Dimension, error) {
	if did <= 0 || uid < 0 {
		return nil, errors.New("did or uid is nil")
	}

	result := &biz.Dimension{}
	if err := d.data.db.Model(&biz.Dimension{}).Where("id = ? and uid = ?", did, uid).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (d *dimensionRepo) GetDmByUidTypeName(ctx context.Context, uid uint, typeStr string, name string) (*biz.Dimension, error) {
	if uid < 0 || typeStr == "" || name == "" {
		return nil, errors.New("uid or type of name is empty")
	}

	result := &biz.Dimension{}
	if err := d.data.db.Model(&biz.Dimension{}).Where("uid = ? and type = ? and name = ?", uid, typeStr, name).First(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (d *dimensionRepo) InsertDm(ctx context.Context, dm *biz.Dimension) error {
	if dm == nil {
		return errors.New("dm is nil")
	}

	return d.data.db.Create(dm).Error
}
func (d *dimensionRepo) UpdateDm(ctx context.Context, dm *biz.Dimension) error {
	if dm == nil {
		return errors.New("dm is nil")
	}

	return d.data.db.Model(&dm).Update("name", dm.Name).Error

}
func (d *dimensionRepo) DeleteDm(ctx context.Context, did uint) error {
	if did <= 0 {
		return errors.New("did is illegal")
	}

	return d.data.db.Delete(&biz.Dimension{}, did).Error
}
func (d *dimensionRepo) GetDmsByType(ctx context.Context, uid uint, typeStr string) ([]*biz.Dimension, error) {
	if uid < 0 || typeStr == "" {
		return nil, errors.New("uid or typeStr is nil")
	}

	result := make([]*biz.Dimension, 0)
	if err := d.data.db.Model(&biz.Dimension{}).Where("uid = ? and type = ?", uid, typeStr).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
