package usecase

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/constants"
	"resource_det_search/internal/utils"
)

type dimensionUsecase struct {
	repo     biz.IDimensionRepo
	userRepo biz.IUserRepo
	cdRepo   biz.IClassDocumentRepo
	docRepo  biz.IDocumentRepo
	logger   *zap.SugaredLogger
}

func NewDimensionUsecase(repo biz.IDimensionRepo, userRepo biz.IUserRepo, cdRepo biz.IClassDocumentRepo, docRepo biz.IDocumentRepo, logger *zap.SugaredLogger) biz.IDimensionUsecase {
	return &dimensionUsecase{
		repo:     repo,
		userRepo: userRepo,
		cdRepo:   cdRepo,
		docRepo:  docRepo,
		logger:   logger,
	}
}

// TODO:维度更新方面涉及到搜索索引更新
// TODO:排序规则后续考虑

func (d *dimensionUsecase) GetUserDm(ctx context.Context, uid uint) (map[string][]*biz.Dimension, error) {
	if uid <= 0 {
		return nil, errors.New("[GetUserDm]the uid is nil")
	}

	// 确认 UID 是否存在
	_, err := d.userRepo.GetUserById(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("[GetUserDm]failed to get GetUserById:err=[%+v]", err)
	}

	dms, err := d.repo.GetDmsByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("[GetUserDm]failed to get GetDmByUid:err=[%+v]", err)
	}

	if len(dms) == 0 {
		return nil, nil
	}

	result := make(map[string][]*biz.Dimension)
	for _, v := range dms {
		if _, ok := result[v.Type]; !ok {
			result[v.Type] = make([]*biz.Dimension, 0)
		}
		result[v.Type] = append(result[v.Type], v)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}
func (d *dimensionUsecase) AddUserDm(ctx context.Context, dm *biz.Dimension) error {
	if dm == nil {
		return errors.New("[AddUserDm]dm is nil")
	}

	reDm, err := d.repo.GetDmByUidTypeName(ctx, dm.Uid, dm.Type, dm.Name)
	if err == nil && reDm.ID > 0 {
		return errors.New("[AddUserDm]the dm is existed")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("[AddUserDm]failed to GetDmByUidTypeName:err=[%+v]", err)
	}

	err = d.repo.InsertDm(ctx, dm)
	if err != nil {
		return fmt.Errorf("[AddUserDm]failed to InsertDm:err=[%+v]", err)
	}

	return nil
}
func (d *dimensionUsecase) UpdateUserDm(ctx context.Context, did uint, name string, uid uint) error {
	if did <= 0 || name == "" || uid <= 0 {
		return errors.New("[UpdateUserDm]did or name or uid is nil")
	}

	// 操作权限（只能用户自身）
	dm, err := d.repo.GetDmById(ctx, did)
	if err != nil {
		return fmt.Errorf("[UpdateUserDm]failed to GetDmById:err=[%+v]", err)
	}

	if dm.Uid != uid {
		return errors.New("[UpdateUserDm]uid difference")
	}

	if dm.Name == name {
		return fmt.Errorf("[UpdateUserDm]update name is same:dm=[%+v]", utils.JsonToString(dm))
	}

	err = d.repo.UpdateDm(ctx, &biz.Dimension{
		Model: gorm.Model{ID: did},
		Name:  name,
	})
	if err != nil {
		return fmt.Errorf("[UpdateUserDm]failed to UpdateDm:err=[%+v]", err)
	}

	// 异步修改搜索引擎存储
	go d.updateDimensionSearch(ctx, did, constants.DmType(dm.Type), dm.Name, name)

	return nil
}
func (d *dimensionUsecase) DeleteUserDm(ctx context.Context, did, uid uint) error {
	if did <= 0 {
		return errors.New("[DeleteUserDm]did is nil")
	}

	// 操作权限（只能用户自身）
	dm, err := d.repo.GetDmById(ctx, did)
	if err != nil {
		return fmt.Errorf("[DeleteUserDm]failed to GetDmById:err=[%+v]", err)
	}

	if dm.Uid != uid {
		return errors.New("[DeleteUserDm]uid difference")
	}

	err = d.repo.DeleteDm(ctx, did)
	if err != nil {
		return fmt.Errorf("[DeleteUserDm]failed to DeleteDm:err=[%+v]", err)
	}

	return nil
}

func (d *dimensionUsecase) GetDmsPartType(ctx context.Context) ([]*biz.Dimension, error) {
	result, err := d.repo.GetDmsPartType(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetDmsPartType]failed to GetDmsPartType:err=[%+v]", err)
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func (d *dimensionUsecase) updateDimensionSearch(ctx context.Context, did uint, typeStr constants.DmType, oldDmName string, newDmName string) {
	defer func() {
		if r := recover(); r != nil {
			d.logger.Errorf("[updateDimensionSearch]panic recover:%+v", r)
		}
	}()

	// select the dm documents docId
	docIds, err := d.docRepo.GetDocIdsByDid(ctx, did)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		d.logger.Errorf("[updateDimensionSearch]failed to GetDocIdsByDid:err=[%+v]", err)
		return
	}

	if len(docIds) == 0 {
		return
	}

	for _, v := range docIds {
		err = d.cdRepo.UpdateDimensions(ctx, v, typeStr, oldDmName, newDmName)
		if err != nil {
			d.logger.Errorf("[updateDimensionSearch]failed to UpdateDimensions:err=[%+v]", err)
		}
	}

	return
}
