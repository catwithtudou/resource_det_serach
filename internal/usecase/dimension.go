package usecase

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"resource_det_search/internal/biz"
)

type dimensionUsecase struct {
	repo     biz.IDimensionRepo
	userRepo biz.IUserRepo
}

func NewDimensionUsecase(repo biz.IDimensionRepo, userRepo biz.IUserRepo) biz.IDimensionUsecase {
	return &dimensionUsecase{
		repo:     repo,
		userRepo: userRepo,
	}
}

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

	err = d.repo.UpdateDm(ctx, &biz.Dimension{
		Model: gorm.Model{ID: did},
		Name:  name,
	})
	if err != nil {
		return fmt.Errorf("[UpdateUserDm]failed to UpdateDm:err=[%+v]", err)
	}

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
