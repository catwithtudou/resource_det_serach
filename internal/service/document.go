package service

import (
	"go.uber.org/zap"
	"resource_det_search/internal/biz"
)

type DocumentService struct {
	log *zap.SugaredLogger
	doc biz.IDocumentUsecase
}

func NewDocumentService(document biz.IDocumentUsecase, logger *zap.SugaredLogger) *DocumentService {
	return &DocumentService{
		log: logger,
		doc: document,
	}
}
