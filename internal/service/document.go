package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"resource_det_search/api"
	v1 "resource_det_search/api/v1"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/constants"
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

func (d *DocumentService) UploadUserDocument(c *gin.Context) {
	uid, _ := c.Get("uid")

	var req v1.UploadUserDocumentReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-UploadUserDocument]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	categories, cgOk := utils.CheckDocTypeStr(req.Categories)
	tags, tgOk := utils.CheckDocTypeStr(req.Tags)
	if len(req.Title) > 50 || len(req.Intro) > 200 || !cgOk || !tgOk {
		d.log.Errorf("[DocumentService-UploadUserDocument]illegal params")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	docFile, err := c.FormFile("doc")
	if err != nil {
		d.log.Errorf("[DocumentService-UploadUserDocument]failed to FormFile:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormFileErr)
		return
	}

	fileType, ok := utils.CheckDocFileType(docFile.Filename)
	if !ok {
		d.log.Errorf("[DocumentService-UploadUserDocument]doc file type not supported:fileName=[%+v]", docFile.Filename)
		c.JSON(http.StatusOK, api.FileTypeErr)
		return
	}

	ok = utils.CheckDocFileSize(docFile.Size)
	if !ok {
		d.log.Errorf("[DocumentService-UploadUserDocument]doc file size not supported:fileSize=[%+v]", docFile.Size)
		c.JSON(http.StatusOK, api.FileSizeErr)
		return
	}

	errCode, err := d.doc.UploadUserDocument(c, &biz.Document{
		Uid:   uid.(uint),
		Type:  fileType,
		Name:  docFile.Filename,
		Intro: req.Intro,
		Title: req.Title,
	}, req.Part, categories, tags, docFile)
	if err != nil {
		d.log.Errorf("[DocumentService-UploadUserDocument]failed to UploadUserDocument:err=[%+v]", err)
		if errCode == constants.DocTitleExist {
			c.JSON(http.StatusOK, api.DocTitleExist)
			return
		}
		if errCode == constants.DocUploadQnyErr {
			c.JSON(http.StatusOK, api.DocUploadQnyErr)
			return
		}
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, api.Success)
	return
}

func (d *DocumentService) GetUserAllDocs(c *gin.Context) {
	var req v1.GetUserAllDocsReq
	var uid uint
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-GetUserAllDocs]failed to bind:err=[%+v]", err)
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

	if !utils.CheckOffsetSize(req.Offset, req.Size) {
		d.log.Errorf("[DocumentService-GetUserAllDocs]illegal params:req=[%+v]", utils.JsonToString(req))
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	docs, docsDmsMap, err := d.doc.GetUserAllDocs(c, uid, req.Offset, req.Size)
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
	resp.Data = make([]*v1.DocPartData, 0, len(docs))
	for _, v := range docs {
		data := &v1.DocPartData{
			DocId:       v.ID,
			Uid:         v.Uid,
			Type:        v.Type,
			Name:        v.Name,
			Intro:       v.Intro,
			Title:       v.Title,
			DownloadNum: v.DownloadNum,
			ScanNum:     v.ScanNum,
			LikeNum:     v.LikeNum,
			UploadTime:  utils.TimeFormat(v.CreatedAt),
			Part:        v1.DimensionUserDmData{},
			Categories:  make([]v1.DimensionUserDmData, 0),
			Tags:        make([]v1.DimensionUserDmData, 0),
		}
		for kk, vv := range docsDmsMap[v.ID] {
			switch kk {
			case string(constants.Part):
				data.Part = v1.DimensionUserDmData{
					Id:   vv[0].ID,
					Name: vv[0].Name,
				}
			case string(constants.Category):
				for _, vvv := range vv {
					data.Categories = append(data.Categories, v1.DimensionUserDmData{
						Id:   vvv.ID,
						Name: vvv.Name,
					})
				}
			case string(constants.Tag):
				for _, vvv := range vv {
					data.Tags = append(data.Tags, v1.DimensionUserDmData{
						Id:   vvv.ID,
						Name: vvv.Name,
					})
				}
			}
		}
		resp.Data = append(resp.Data, data)
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (d *DocumentService) GetAllDocs(c *gin.Context) {
	var req v1.GetAllDocsReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-GetAllDocs]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if !utils.CheckOffsetSize(req.Offset, req.Size) || req.Offset < 0 || !utils.CheckSortBy(req.SortBy) {
		d.log.Errorf("[DocumentService-GetAllDocs]illegal params:req=[%+v]", utils.JsonToString(req.SortBy))
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	docs, docsDmsMap, err := d.doc.GetAllDocs(c, req.Offset, req.Size, req.SortBy)
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
	resp.Data = make([]*v1.DocPartData, 0, len(docs))
	for _, v := range docs {
		data := &v1.DocPartData{
			DocId:       v.ID,
			Uid:         v.Uid,
			Type:        v.Type,
			Name:        v.Name,
			Intro:       v.Intro,
			Title:       v.Title,
			DownloadNum: v.DownloadNum,
			ScanNum:     v.ScanNum,
			LikeNum:     v.LikeNum,
			UploadTime:  utils.TimeFormat(v.CreatedAt),
			Part:        v1.DimensionUserDmData{},
			Categories:  make([]v1.DimensionUserDmData, 0),
			Tags:        make([]v1.DimensionUserDmData, 0),
		}
		for kk, vv := range docsDmsMap[v.ID] {
			switch kk {
			case string(constants.Part):
				data.Part = v1.DimensionUserDmData{
					Id:   vv[0].ID,
					Name: vv[0].Name,
				}
			case string(constants.Category):
				for _, vvv := range vv {
					data.Categories = append(data.Categories, v1.DimensionUserDmData{
						Id:   vvv.ID,
						Name: vvv.Name,
					})
				}
			case string(constants.Tag):
				for _, vvv := range vv {
					data.Tags = append(data.Tags, v1.DimensionUserDmData{
						Id:   vvv.ID,
						Name: vvv.Name,
					})
				}
			}
		}
		resp.Data = append(resp.Data, data)
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

	if req.Did <= 0 || !utils.CheckOffsetSize(req.Offset, req.Size) {
		d.log.Errorf("[DocumentService-GetUserDimensionDocs]illegal params:req=[%+v]", utils.JsonToString(req))
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	docs, dm, err := d.doc.GetDmDocs(c, uid, req.Did, req.Offset, req.Size)
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

	resp.Data.Docs = make([]*v1.DocNoDmData, 0, len(docs))
	for _, v := range docs {
		resp.Data.Docs = append(resp.Data.Docs, &v1.DocNoDmData{
			DocId:       v.ID,
			Uid:         v.Uid,
			Type:        v.Type,
			Name:        v.Name,
			Intro:       v.Intro,
			Title:       v.Title,
			DownloadNum: v.DownloadNum,
			ScanNum:     v.ScanNum,
			LikeNum:     v.LikeNum,
			UploadTime:  utils.TimeFormat(v.CreatedAt),
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

	if !utils.CheckUserType(req.Type) || !utils.CheckOffsetSize(req.Offset, req.Size) {
		d.log.Errorf("[DocumentService-GetUserAllDimensionDocs]illegal params:req=[%+v]", utils.JsonToString(req))
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	typeDocs, err := d.doc.GetAllDmTypeDocs(c, uid, req.Type, req.Offset, req.Size)
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

	if req.Did <= 0 || !utils.CheckOffsetSize(req.Offset, req.Size) || !utils.CheckSortBy(req.SortBy) {
		d.log.Errorf("[DocumentService-GetDimensionDocs]illegal params:req=[%+v]", utils.JsonToString(req))
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	docs, docsDmsMap, err := d.doc.GetPartDocs(c, req.Did, req.Offset, req.Size, req.SortBy)
	if err != nil {
		d.log.Errorf("[DocumentService-GetDimensionDocs]failed to GetDmDocs:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.GetDimensionDocsResp{
		RespCommon: api.Success,
	}
	if len(docs) == 0 {
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = make([]*v1.DocPartData, 0, len(docs))
	for _, v := range docs {
		data := &v1.DocPartData{
			DocId:       v.ID,
			Uid:         v.Uid,
			Type:        v.Type,
			Name:        v.Name,
			Intro:       v.Intro,
			Title:       v.Title,
			DownloadNum: v.DownloadNum,
			ScanNum:     v.ScanNum,
			LikeNum:     v.LikeNum,
			UploadTime:  utils.TimeFormat(v.CreatedAt),
			Part:        v1.DimensionUserDmData{},
			Categories:  make([]v1.DimensionUserDmData, 0),
			Tags:        make([]v1.DimensionUserDmData, 0),
		}
		for kk, vv := range docsDmsMap[v.ID] {
			switch kk {
			case string(constants.Part):
				data.Part = v1.DimensionUserDmData{
					Id:   vv[0].ID,
					Name: vv[0].Name,
				}
			case string(constants.Category):
				for _, vvv := range vv {
					data.Categories = append(data.Categories, v1.DimensionUserDmData{
						Id:   vvv.ID,
						Name: vvv.Name,
					})
				}
			case string(constants.Tag):
				for _, vvv := range vv {
					data.Tags = append(data.Tags, v1.DimensionUserDmData{
						Id:   vvv.ID,
						Name: vvv.Name,
					})
				}
			}
		}
		resp.Data = append(resp.Data, data)
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

	if !utils.CheckAllType(req.Type) || !utils.CheckOffsetSize(req.Offset, req.Size) {
		d.log.Errorf("[DocumentService-GetAllDimensionDocs]illegal params:req=[%+v]", req)
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	typeDocs, err := d.doc.GetAllDmTypeDocs(c, 0, req.Type, req.Offset, req.Size)
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

func (d *DocumentService) DetUserDoc(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		d.log.Errorf("[DocumentService-DetUserDoc]failed to FormFile:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormFileErr)
		return
	}

	fileType, ok := utils.CheckDocFileType(file.Filename)
	if !ok {
		d.log.Errorf("[DocumentService-DetUserDoc]failed to FormFile:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FileTypeErr)
		return
	}

	detail, err := d.doc.DetFile(c, fileType, file)
	if err != nil {
		d.log.Errorf("[DocumentService-DetUserDoc]failed to DetFile:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, v1.DetUserDocResp{
		RespCommon: api.Success,
		Data:       detail,
	})
	return
}

func (d *DocumentService) GetDocWithDms(c *gin.Context) {
	var req v1.GetDocWithDmsReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-GetDocWithDms]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	doc, dmMap, err := d.doc.GetDocWithDms(c, req.DocId)
	if err != nil {
		d.log.Errorf("[DocumentService-GetDocWithDms]failed to GetDocWithDms:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.GetDocWithDmsResp{
		RespCommon: api.Success,
		Data: &v1.DocWithDmsData{
			DocId:       doc.ID,
			Uid:         doc.Uid,
			Type:        doc.Type,
			Dir:         doc.Dir,
			Name:        doc.Name,
			Intro:       doc.Intro,
			Title:       doc.Title,
			DownloadNum: doc.DownloadNum,
			ScanNum:     doc.ScanNum,
			LikeNum:     doc.LikeNum,
			Content:     doc.Content,
			Part:        v1.DimensionUserDmData{},
			Categories:  make([]v1.DimensionUserDmData, 0),
			Tags:        make([]v1.DimensionUserDmData, 0),
		},
	}
	for k, v := range dmMap {
		switch k {
		case string(constants.Part):
			resp.Data.Part = v1.DimensionUserDmData{
				Id:   v[0].ID,
				Name: v[0].Name,
			}
		case string(constants.Category):
			for _, vv := range v {
				resp.Data.Categories = append(resp.Data.Categories, v1.DimensionUserDmData{
					Id:   vv.ID,
					Name: vv.Name,
				})
			}
		case string(constants.Tag):
			for _, vv := range v {
				resp.Data.Tags = append(resp.Data.Tags, v1.DimensionUserDmData{
					Id:   vv.ID,
					Name: vv.Name,
				})
			}
		}
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (d *DocumentService) GetUserDocCount(c *gin.Context) {
	var req v1.GetUserDocCountReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DocumentService-GetUserDocCount]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	count, err := d.doc.GetUserDocCount(c, req.Uid)
	if err != nil {
		d.log.Errorf("[DocumentService-GetUserDocCount]failed to GetUserAllDocs:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, &v1.GetUserDocCountResp{
		RespCommon: api.Success,
		Data:       &v1.UserDocCountData{Count: count},
	})
	return
}
