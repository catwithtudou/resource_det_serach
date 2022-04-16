package usecase

import (
	"context"
	"errors"
	"fmt"
	"resource_det_search/internal/biz"
)

type documentUsecase struct {
	repo     biz.IDocumentRepo
	userRepo biz.IUserRepo
	dmRepo   biz.IDimensionRepo
}

func NewDocumentUsecase(repo biz.IDocumentRepo, userRepo biz.IUserRepo, dmRepo biz.IDimensionRepo) biz.IDocumentUsecase {
	return &documentUsecase{
		repo:     repo,
		userRepo: userRepo,
		dmRepo:   dmRepo,
	}
}

func (d *documentUsecase) GetUserAllDocs(ctx context.Context, uid uint) ([]*biz.Document, error) {
	if uid <= 0 {
		return nil, errors.New("[GetUserAllDocs]the uid is nil")
	}

	res, err := d.repo.GetDocsByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("[GetUserAllDocs]failed to GetDocsByUid:err=[%+v]", err)
	}
	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
}
func (d *documentUsecase) GetAllDocs(ctx context.Context) ([]*biz.Document, error) {

	res, err := d.repo.GetDocs(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetAllDocs]failed to GetDocs:err=[%+v]", err)
	}
	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
}
func (d *documentUsecase) GetDmDocs(ctx context.Context, uid uint, did uint) ([]*biz.Document, *biz.Dimension, error) {
	if uid <= 0 || did <= 0 {
		return nil, nil, errors.New("[GetDmDocs]uid or did is nil")
	}

	// select the dimension info
	dmInfo, err := d.dmRepo.GetDmByDidUid(ctx, did, uid)
	if err != nil {
		return nil, nil, fmt.Errorf("[GetDmDocs]failed to GetDmById:err=[%+v]", err)
	}
	if dmInfo == nil {
		return nil, nil, errors.New("[GetDmDocs]dmInfo is nil")
	}

	// select the docs with did
	docs, err := d.repo.GetDocsWithDid(ctx, did)
	if err != nil {
		return nil, nil, fmt.Errorf("[GetDmDocs]failed to GetDocsWithDid:err=[%+v]", err)
	}
	if len(docs) == 0 {
		return nil, dmInfo, nil
	}

	return docs, dmInfo, nil
}
func (d *documentUsecase) GetAllDmTypeDocs(ctx context.Context, uid uint, typeStr string) (map[string][]*biz.Document, error) {
	if uid <= 0 || typeStr == "" {
		return nil, errors.New("[GetAllDmTypeDocs]uid or typeStr is nil")
	}

	// select the user dimension infos with dimension type
	dmInfos, err := d.dmRepo.GetDmsByType(ctx, uid, typeStr)
	if err != nil {
		return nil, fmt.Errorf("[GetAllDmTypeDocs]failed to GetDmsByType:err=[%+v]", err)
	}
	if len(dmInfos) == 0 {
		return nil, nil
	}

	// for select the docs with the dm id
	result := make(map[string][]*biz.Document)
	for _, v := range dmInfos {
		result[v.Name] = make([]*biz.Document, 0)
		docs, err := d.repo.GetDocsWithDid(ctx, v.ID)
		if err != nil {
			return nil, fmt.Errorf("[GetAllDmTypeDocs]failed to GetDocsWithDid:err=[%+v]", err)
		}
		result[v.Name] = docs
	}

	return result, nil
}
func (d *documentUsecase) AddLikeDoc(ctx context.Context, docId uint, num uint) error {
	if docId <= 0 || num <= 0 {
		return errors.New("[AddLikeDoc]docId or num is nil")
	}

	if err := d.repo.AddDocLikeNum(ctx, docId, num); err != nil {
		return fmt.Errorf("[AddLikeDoc]failed to AddDocLikeNum:err=[%+v]", err)
	}

	return nil
}

func (d *documentUsecase) DeleteUserDoc(ctx context.Context, docId uint, uid uint) error {
	if docId <= 0 || uid <= 0 {
		return errors.New("[DeleteUserDoc]docId or uid is nil")
	}

	if err := d.repo.DeleteDocByIdWithUid(ctx, docId, uid); err != nil {
		return fmt.Errorf("[DeleteUserDoc]failed to DeleteDocByIdWithUid:err=[%+v]", err)
	}
	return nil
}
