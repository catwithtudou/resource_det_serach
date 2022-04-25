package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"resource_det_search/api"
	v1 "resource_det_search/api/v1"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/utils"
)

type ClassDocumentService struct {
	log *zap.SugaredLogger
	cd  biz.IClassDocumentUsecase
}

func NewClassDocumentService(cd biz.IClassDocumentUsecase, logger *zap.SugaredLogger) *ClassDocumentService {
	return &ClassDocumentService{
		log: logger,
		cd:  cd,
	}
}

func (cd *ClassDocumentService) SearchAllQuery(c *gin.Context) {
	var req v1.SearchAllQueryReq
	if err := c.ShouldBind(&req); err != nil {
		cd.log.Errorf("[ClassDocumentService-SearchAllQuery]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if !utils.CheckSearchSortBy(req.SortBy) {
		cd.log.Errorf("[ClassDocumentService-SearchAllQuery]illegal params:req=[%+v]", utils.JsonToString(req))
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	cds, err := cd.cd.SearchAllQuery(c, req.Detail, req.PartId, req.SortBy)
	if err != nil {
		cd.log.Errorf("[ClassDocumentService-SearchAllQuery]failed to SearchAllQuery:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := v1.SearchAllQueryResp{
		RespCommon: api.Success,
		Data:       make([]*v1.ClassDocumentData, len(cds)),
	}
	for i := 0; i < len(cds); i++ {
		resp.Data[i] = &v1.ClassDocumentData{
			DocId:       cds[i].Id,
			Title:       cds[i].Title,
			Content:     cds[i].Content,
			Intro:       cds[i].Intro,
			Part:        cds[i].Part,
			Categories:  make([]string, 0),
			Tags:        make([]string, 0),
			FileType:    cds[i].FileType,
			Username:    cds[i].Username,
			UploadTime:  utils.TimestampFormat(cds[i].UploadDate),
			DownloadNum: cds[i].DownloadNum,
			ScanNum:     cds[i].ScanNum,
			LikeNum:     cds[i].LikeNum,
		}
		if len(cds[i].Categories) > 0 {
			resp.Data[i].Categories = cds[i].Categories
		}
		if len(cds[i].Tags) > 0 {
			resp.Data[i].Tags = cds[i].Tags
		}
	}
	c.JSON(http.StatusOK, resp)
	return
}
