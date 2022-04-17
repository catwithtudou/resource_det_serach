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

func (d *DocumentService) GetUserAllDocs(c *gin.Context) {
	var req v1.DimensionGetUserDmReq
	var uid uint
	if _ = c.ShouldBind(&req); req.Uid > 0 {
		uid = req.Uid
	}

	if uid <= 0 {
		getUid, _ := c.Get("uid")
		uid = getUid.(uint)
	}

	docs, err := d.doc.GetUserAllDocs(c, uid)
	if err != nil {
		d.log.Errorf("[DocumentService-GetUserAllDocs]failed to GetUserAllDocs:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.GetUserAllDocsResp{
		RespCommon: api.Success,
	}
	if len(docs) == 0 {
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = make([]*v1.DocData, 0, len(docs))
	for _, v := range docs {
		resp.Data = append(resp.Data, &v1.DocData{
			DocId:       v.ID,
			Uid:         v.Uid,
			Type:        v.Type,
			Dir:         v.Dir,
			Name:        v.Name,
			Intro:       v.Intro,
			Title:       v.Title,
			DownloadNum: v.DownloadNum,
			ScanNum:     v.ScanNum,
			LikeNum:     v.LikeNum,
			Content:     v.Content,
		})
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (d *DocumentService) GetAllDocs(c *gin.Context) {
	docs, err := d.doc.GetAllDocs(c)
	if err != nil {
		d.log.Errorf("[DocumentService-GetAllDocs]failed to GetUserAllDocs:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.GetAllDocsResp{
		RespCommon: api.Success,
	}
	if len(docs) == 0 {
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = make([]*v1.DocData, 0, len(docs))
	for _, v := range docs {
		resp.Data = append(resp.Data, &v1.DocData{
			DocId:       v.ID,
			Uid:         v.Uid,
			Type:        v.Type,
			Dir:         v.Dir,
			Name:        v.Name,
			Intro:       v.Intro,
			Title:       v.Title,
			DownloadNum: v.DownloadNum,
			ScanNum:     v.ScanNum,
			LikeNum:     v.LikeNum,
			Content:     v.Content,
		})
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (d *DocumentService) GetUserDimensionDocs(c *gin.Context) {
	var req v1.GetUserDimensionDocsReq
	var uid uint
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-GetUserDimensionDocs]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if req.Uid > 0 {
		uid = req.Uid
	}

	if uid <= 0 {
		getUid, _ := c.Get("uid")
		uid = getUid.(uint)
	}

	if req.Did <= 0 {
		d.log.Errorf("[DocumentService-GetUserDimensionDocs]illegal did")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	docs, dm, err := d.doc.GetDmDocs(c, uid, req.Did)
	if err != nil {
		d.log.Errorf("[DocumentService-GetUserDimensionDocs]failed to GetDmDocs:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.GetUserDimensionDocsResp{
		RespCommon: api.Success,
		Data: &v1.DimensionDocsData{
			Did:    dm.ID,
			DmName: dm.Name,
			DmType: dm.Type,
		},
	}
	if len(docs) == 0 {
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data.Docs = make([]*v1.DocData, 0, len(docs))
	for _, v := range docs {
		resp.Data.Docs = append(resp.Data.Docs, &v1.DocData{
			DocId:       v.ID,
			Uid:         v.Uid,
			Type:        v.Type,
			Dir:         v.Dir,
			Name:        v.Name,
			Intro:       v.Intro,
			Title:       v.Title,
			DownloadNum: v.DownloadNum,
			ScanNum:     v.ScanNum,
			LikeNum:     v.LikeNum,
			Content:     v.Content,
		})
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (d *DocumentService) GetUserAllDimensionDocs(c *gin.Context) {
	var req v1.GetUserAllDimensionDocsReq
	var uid uint
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-GetUserAllDimensionDocs]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if req.Uid > 0 {
		uid = req.Uid
	}

	if uid <= 0 {
		getUid, _ := c.Get("uid")
		uid = getUid.(uint)
	}

	if !utils.CheckUserType(req.Type) {
		d.log.Errorf("[DocumentService-GetUserAllDimensionDocs]illegal params")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	typeDocs, err := d.doc.GetAllDmTypeDocs(c, uid, req.Type)
	if err != nil {
		d.log.Errorf("[DocumentService-GetUserAllDimensionDocs]failed to GetAllDmTypeDocs:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.GetUserAllDimensionDocsResp{
		RespCommon: api.Success,
		Data:       make(map[string][]*v1.DocData),
	}

	for k, v := range typeDocs {
		if _, ok := resp.Data[k]; !ok {
			resp.Data[k] = make([]*v1.DocData, 0, len(v))
		}
		for _, vv := range v {
			resp.Data[k] = append(resp.Data[k], &v1.DocData{
				DocId:       vv.ID,
				Uid:         vv.Uid,
				Type:        vv.Type,
				Dir:         vv.Dir,
				Name:        vv.Name,
				Intro:       vv.Intro,
				Title:       vv.Title,
				DownloadNum: vv.DownloadNum,
				ScanNum:     vv.ScanNum,
				LikeNum:     vv.LikeNum,
				Content:     vv.Content,
			})
		}
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (d *DocumentService) GetDimensionDocs(c *gin.Context) {
	var req v1.GetDimensionDocsReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-GetDimensionDocs]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if req.Did <= 0 {
		d.log.Errorf("[DocumentService-GetDimensionDocs]illegal did")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	docs, dm, err := d.doc.GetDmDocs(c, 0, req.Did)
	if err != nil {
		d.log.Errorf("[DocumentService-GetUserDimensionDocs]failed to GetDmDocs:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.GetDimensionDocsResp{
		RespCommon: api.Success,
		Data: &v1.DimensionDocsData{
			Did:    dm.ID,
			DmName: dm.Name,
			DmType: dm.Type,
		},
	}
	if len(docs) == 0 {
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data.Docs = make([]*v1.DocData, 0, len(docs))
	for _, v := range docs {
		resp.Data.Docs = append(resp.Data.Docs, &v1.DocData{
			DocId:       v.ID,
			Uid:         v.Uid,
			Type:        v.Type,
			Dir:         v.Dir,
			Name:        v.Name,
			Intro:       v.Intro,
			Title:       v.Title,
			DownloadNum: v.DownloadNum,
			ScanNum:     v.ScanNum,
			LikeNum:     v.LikeNum,
			Content:     v.Content,
		})
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (d *DocumentService) GetAllDimensionDocs(c *gin.Context) {
	var req v1.GetAllDimensionDocsReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-GetAllDimensionDocs]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if !utils.CheckAllType(req.Type) {
		d.log.Errorf("[DocumentService-GetAllDimensionDocs]illegal params")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	typeDocs, err := d.doc.GetAllDmTypeDocs(c, 0, req.Type)
	if err != nil {
		d.log.Errorf("[DocumentService-GetAllDimensionDocs]failed to GetAllDmTypeDocs:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.GetAllDimensionDocsResp{
		RespCommon: api.Success,
		Data:       make(map[string][]*v1.DocData),
	}

	for k, v := range typeDocs {
		if _, ok := resp.Data[k]; !ok {
			resp.Data[k] = make([]*v1.DocData, 0, len(v))
		}
		for _, vv := range v {
			resp.Data[k] = append(resp.Data[k], &v1.DocData{
				DocId:       vv.ID,
				Uid:         vv.Uid,
				Type:        vv.Type,
				Dir:         vv.Dir,
				Name:        vv.Name,
				Intro:       vv.Intro,
				Title:       vv.Title,
				DownloadNum: vv.DownloadNum,
				ScanNum:     vv.ScanNum,
				LikeNum:     vv.LikeNum,
				Content:     vv.Content,
			})
		}
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (d *DocumentService) AddLikeDoc(c *gin.Context) {
	var req v1.AddLikeDocReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-AddLikeDoc]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	err := d.doc.AddLikeDoc(c, req.DocId, 1)
	if err != nil {
		d.log.Errorf("[DocumentService-AddLikeDoc]failed to AddLikeDoc:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, api.Success)
	return
}

func (d *DocumentService) DeleteUserDoc(c *gin.Context) {
	uid, _ := c.Get("uid")

	var req v1.DeleteUserDocReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-DeleteUserDoc]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	err := d.doc.DeleteUserDoc(c, req.DocId, uid.(uint))
	if err != nil {
		d.log.Errorf("[DocumentService-DeleteUserDoc]failed to DeleteUserDoc:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, api.Success)
	return
}
