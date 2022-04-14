package usecase

import (
	"context"
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
	return nil, nil
}
func (d *documentUsecase) GetAllDocs(ctx context.Context) ([]*biz.Document, error) {
	return nil, nil
}
func (d *documentUsecase) GetDmDocs(ctx context.Context, uid uint, did uint) ([]*biz.Document, string, uint, error) {
	return nil, "", 0, nil
}
func (d *documentUsecase) GetAllDmTypeDocs(ctx context.Context, uid uint, typeStr string) (map[string][]*biz.Document, string, error) {
	return nil, "", nil
}
func (d *documentUsecase) AddLikeDoc(ctx context.Context, docId uint) error {
	return nil
}
func (d *documentUsecase) DeleteUserDoc(ctx context.Context, docId uint, uid uint) error {
	return nil
}
