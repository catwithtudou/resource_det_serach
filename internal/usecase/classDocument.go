package usecase

import (
	"context"
	"errors"
	"fmt"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/constants"
	"resource_det_search/internal/utils"
)

type classDocumentUsecase struct {
	repo biz.IClassDocumentRepo
	dm   biz.IDimensionRepo
}

func NewClassDocumentUsecase(repo biz.IClassDocumentRepo, dm biz.IDimensionRepo) biz.IClassDocumentUsecase {
	return &classDocumentUsecase{
		repo: repo,
		dm:   dm,
	}
}

func (c *classDocumentUsecase) SearchAllQuery(ctx context.Context, queryStr string, partId uint, sortBy string) ([]*biz.ClassDocument, error) {
	if queryStr == "" {
		return nil, errors.New("[SearchAllQuery]queryStr is empty")
	}

	if partId <= 0 {
		res, err := c.repo.SearchAllQuery(ctx, queryStr, sortBy)
		if err != nil {
			return nil, fmt.Errorf("[SearchAllQuery]failed to SearchAllQuery:err=[%+v]", err)
		}
		return res, err
	}

	dmInfo, err := c.dm.GetDmById(ctx, partId)
	if err != nil {
		return nil, fmt.Errorf("[SearchAllQuery]failed to GetDmById:er=[%+v]", err)
	}
	if dmInfo.Type != string(constants.Part) {
		return nil, fmt.Errorf("[SearchAllQuery]did is illeagal:dm=[%+v]", utils.JsonToString(dmInfo))
	}

	res, err := c.repo.SearchQueryByPart(ctx, queryStr, dmInfo.Name, sortBy)
	if err != nil {
		return nil, fmt.Errorf("[SearchAllQuery]failed to SearchQueryByPart:er=[%+v]", err)
	}
	return res, nil
}
