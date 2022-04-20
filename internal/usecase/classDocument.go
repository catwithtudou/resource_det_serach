package usecase

import (
	"context"
	"errors"
	"fmt"
	"resource_det_search/internal/biz"
)

type classDocumentUsecase struct {
	repo biz.IClassDocumentRepo
}

func NewClassDocumentUsecase(repo biz.IClassDocumentRepo) biz.IClassDocumentUsecase {
	return &classDocumentUsecase{
		repo: repo,
	}
}

func (c *classDocumentUsecase) SearchAllQuery(ctx context.Context, queryStr string) ([]*biz.ClassDocument, error) {
	if queryStr == "" {
		return nil, errors.New("[SearchAllQuery]queryStr is empty")
	}

	res, err := c.repo.SearchAllQuery(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("[SearchAllQuery]failed to SearchAllQuery:err=[%+v]", err)
	}

	return res, nil
}
