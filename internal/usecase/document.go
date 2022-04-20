package usecase

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/constants"
	"resource_det_search/internal/utils"
	"time"
)

type documentUsecase struct {
	repo     biz.IDocumentRepo
	userRepo biz.IUserRepo
	dmRepo   biz.IDimensionRepo
	cdRepo   biz.IClassDocumentRepo
	logger   *zap.SugaredLogger
}

func NewDocumentUsecase(repo biz.IDocumentRepo, userRepo biz.IUserRepo, dmRepo biz.IDimensionRepo, cdRepo biz.IClassDocumentRepo, logger *zap.SugaredLogger) biz.IDocumentUsecase {
	return &documentUsecase{
		repo:     repo,
		userRepo: userRepo,
		dmRepo:   dmRepo,
		cdRepo:   cdRepo,
		logger:   logger,
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

	result := make([]*biz.Document, 0, len(res))
	for _, v := range result {
		// if document is_load_search and  is_sava are false,it should be continued
		if !v.IsLoadSearch || !v.IsSave {
			continue
		}
		result = append(result, v)
	}

	return result, nil
}
func (d *documentUsecase) GetAllDocs(ctx context.Context) ([]*biz.Document, error) {

	res, err := d.repo.GetDocs(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetAllDocs]failed to GetDocs:err=[%+v]", err)
	}
	if len(res) == 0 {
		return nil, nil
	}

	result := make([]*biz.Document, 0, len(res))
	for _, v := range result {
		// if document is_load_search and  is_sava are false,it should be continued
		if !v.IsLoadSearch || !v.IsSave {
			continue
		}
		result = append(result, v)
	}

	return result, nil
}
func (d *documentUsecase) GetDmDocs(ctx context.Context, uid uint, did uint) ([]*biz.Document, *biz.Dimension, error) {
	if uid < 0 || did <= 0 {
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

	reDocs := make([]*biz.Document, 0, len(docs))
	for _, v := range reDocs {
		// if document is_load_search and  is_sava are false,it should be continued
		if !v.IsLoadSearch || !v.IsSave {
			continue
		}
		reDocs = append(reDocs, v)
	}

	return reDocs, dmInfo, nil
}
func (d *documentUsecase) GetAllDmTypeDocs(ctx context.Context, uid uint, typeStr string) (map[string][]*biz.Document, error) {
	if uid < 0 || typeStr == "" {
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
		for _, vv := range docs {
			// if document is_load_search and  is_sava are false,it should be continued
			if !vv.IsLoadSearch || !vv.IsSave {
				continue
			}
			result[v.Name] = append(result[v.Name], vv)
		}
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

	if err := d.repo.DeleteDocWithDmsByIdWithUid(ctx, docId, uid); err != nil {
		return fmt.Errorf("[DeleteUserDoc]failed to DeleteDocByIdWithUid:err=[%+v]", err)
	}
	return nil
}
func (d *documentUsecase) UploadUserDocument(ctx context.Context, doc *biz.Document, part uint, categories []uint, tags []uint, fileData *multipart.FileHeader) (constants.ErrCode, error) {
	if doc == nil || part <= 0 || fileData == nil {
		return constants.DefaultErr, errors.New("[UploadUserDocument]doc or dmIds or fileData is nil")
	}

	// select the dmIds illegal
	partDm, err := d.dmRepo.GetDmById(ctx, part)
	if err != nil {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]failed to GetDmById:err=[%+v]", err)
	}
	if partDm == nil || partDm.Type != string(constants.Part) {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]illegal part id:dm=[%+v]", utils.JsonToString(partDm))
	}
	categoriesDm, err := d.checkDmIdsIllegal(ctx, doc.Uid, constants.Category, categories)
	if err != nil {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]illegal category id:err=[%+v],ids=[%+v]", err, categories)
	}
	tagsDm, err := d.checkDmIdsIllegal(ctx, doc.Uid, constants.Tag, tags)
	if err != nil {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]illegal tag id:err=[%+v],ids=[%+v]", err, tags)
	}

	// select the repo to judge repeat and insert the repo
	err = d.repo.GetSaveDocWithNameAndTitle(ctx, doc.Uid, doc.Title)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]failed to GetSaveDocWithNameAndTitle:err=[%+v]", err)
	}
	if err == nil {
		return constants.DocTitleExist, errors.New("[UploadUserDocument]illegal title")
	}

	// insert the repo
	docId, err := d.repo.InsertDocWithDms(ctx, doc, d.allDmTypeIdsToIds(part, categories, tags))
	if err != nil {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]failed to InsertDocWithDms:err=[%+v]", err)
	}

	// upload the qny file
	fileBytes, err := utils.MultipartFileHeaderToBytes(fileData)
	if err != nil {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]failed to MultipartFileHeaderToBytes:err=[%+v]", err)
	}

	key, err := utils.UploadPartByteData(ctx, fileBytes, utils.GenDocKey(docId, doc.Uid))
	if err != nil {
		return constants.DocUploadQnyErr, fmt.Errorf("[UploadUserDocument]failed to UploadByteData:err=[%+v]", err)
	}

	// handle done update the repo(todo:async goroutine handle)
	err = d.repo.UpdateDocById(ctx, &biz.Document{
		Model:  gorm.Model{ID: docId},
		Dir:    utils.GenFileLink(key),
		IsSave: true,
	})
	if err != nil {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]failed to UpdateDocById:err=[%+v]", err)
	}

	// async goroutine handle:det file details (with update the database) and upload the es search
	go d.uploadDetSearch(ctx, docId, doc, partDm, categoriesDm, tagsDm, fileBytes)

	return constants.Success, nil
}
func (d *documentUsecase) DetFile(ctx context.Context, fileType string, fileData *multipart.FileHeader) (string, error) {
	if fileData == nil {
		return "", errors.New("[DetFile]fileData is nil")
	}

	fileBytes, err := utils.MultipartFileHeaderToBytes(fileData)
	if err != nil {
		return "", fmt.Errorf("[DetFile]failed to MultipartFileHeaderToBytes:err=[%+v]", err)
	}

	res, err := d.detFile(fileType, fileBytes)
	if err != nil {
		return "", fmt.Errorf("[DetFile]failed to detFile:err=[%+v],fileType=[%+v]", err, fileType)
	}

	return res, nil
}
func (d *documentUsecase) GetDocWithDms(ctx context.Context, docId uint) (*biz.Document, map[string][]*biz.Dimension, error) {
	if docId <= 0 {
		return nil, nil, errors.New("[GetDocWithDms]docId is nil")
	}

	doc, dms, err := d.repo.GetDocWithDms(ctx, docId)
	if err != nil {
		return nil, nil, fmt.Errorf("[GetDocWithDms]failed to GetDocWithDms:err=[%+v],docId=[%+v]", err, docId)
	}
	if doc == nil || len(dms) == 0 {
		return nil, nil, fmt.Errorf("[GetDocWithDms]doc or dms is nil:docId=[%+v]", docId)
	}

	dmMap := make(map[string][]*biz.Dimension)
	for _, v := range dms {
		if _, ok := dmMap[v.Type]; !ok {
			dmMap[v.Type] = make([]*biz.Dimension, 0)
		}
		dmMap[v.Type] = append(dmMap[v.Type], v)
	}

	//todo:增加浏览量+下载量

	return doc, dmMap, nil
}

func (d *documentUsecase) uploadDetSearch(ctx context.Context, docId uint, doc *biz.Document, part *biz.Dimension, categories []*biz.Dimension, tags []*biz.Dimension, fileData []byte) {
	defer func() {
		if r := recover(); r != nil {
			d.logger.Errorf("[uploadDetSearch]panic recover:%+v", r)
		}
	}()

	// det file
	// todo:识别失败将content作为空处理
	res, err := d.detFile(doc.Type, fileData)
	if err != nil {
		d.logger.Errorf("[uploadDetSearch]failed to detFile:err=[%+v],doc=[%+v]", err, utils.JsonToString(doc))
		return
	}

	// upload search file
	user, err := d.userRepo.GetUserById(ctx, doc.Uid)
	if err != nil {
		d.logger.Errorf("[uploadDetSearch]failed to GetUserById:err=[%+v],doc=[%+v]", err, utils.JsonToString(doc))
		return
	}
	cd := &biz.ClassDocument{
		Id:         docId,
		Title:      doc.Title,
		Content:    res,
		Intro:      doc.Intro,
		Part:       part.Name,
		FileType:   doc.Type,
		Username:   user.Name,
		UploadDate: time.Now().Unix(),
	}
	if len(categories) > 0 {
		cd.Categories = make([]string, 0, len(categories))
		for _, v := range categories {
			cd.Categories = append(cd.Categories, v.Name)
		}
	}
	if len(tags) > 0 {
		cd.Tags = make([]string, 0, len(tags))
		for _, v := range tags {
			cd.Tags = append(cd.Tags, v.Name)
		}
	}

	err = d.cdRepo.InsertDoc(ctx, docId, cd)
	if err != nil {
		d.logger.Errorf("[uploadDetSearch]failed to InsertDoc:err=[%+v],doc=[%+v],cd=[%+v]", err, utils.JsonToString(doc), utils.JsonToString(cd))
		return
	}

	// update docRepo
	err = d.repo.UpdateDocById(ctx, &biz.Document{
		Model:        gorm.Model{ID: docId},
		IsLoadSearch: true,
		Content:      res,
	})
	if err != nil {
		d.logger.Errorf("[uploadDetSearch]failed to UpdateDocById:err=[%+v],doc=[%+v],cd=[%+v]", err, utils.JsonToString(doc), utils.JsonToString(cd))
		return
	}

	return
}

func (d *documentUsecase) checkDmIdsIllegal(ctx context.Context, uid uint, typeStr constants.DmType, ids []uint) ([]*biz.Dimension, error) {
	if len(ids) <= 0 {
		return nil, nil
	}

	result, err := d.dmRepo.GetDmsInIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		if v.Uid != uid || v.Type != string(typeStr) {
			return nil, errors.New("illegal dm id")
		}
	}

	return result, nil
}

func (d *documentUsecase) allDmTypeIdsToIds(part uint, categories []uint, tags []uint) []uint {
	result := make([]uint, 0, len(categories)+len(tags)+1)
	result = append(result, part)
	if len(categories) > 0 {
		result = append(result, categories...)
	}
	if len(tags) > 0 {
		result = append(result, tags...)
	}

	return result
}

func (d *documentUsecase) detFile(fileType string, fileBytes []byte) (string, error) {
	if utils.DetByteTypesContains(fileType) {
		detType := constants.DetByteType(fileType)
		switch detType {
		case constants.Txt:
			return string(fileBytes), nil
		case constants.Docx:
			txt, err := utils.DetDocxByUnidoc(fileBytes)
			if err != nil {
				return "", err
			}
			return txt, nil
		case constants.Pptx:
			txt, err := utils.DetPptxByUnidoc(fileBytes)
			if err != nil {
				return "", err
			}
			return txt, nil
		case constants.Xlsx:
			txt, err := utils.DetXlsxByUnidoc(fileBytes)
			if err != nil {
				return "", err
			}
			return txt, nil
		case constants.Md:
			return utils.DetMd(fileBytes)
		default:
			return "", errors.New("det byte type not supported")
		}

	}

	// OCR识别部分
	if utils.DetOcrTypesContains(fileType) {
		detType := constants.DetOcrType(fileType)
		switch detType {
		default:
			return "", errors.New("det ocr type not supported")
		}

	}

	return "", errors.New("det doc type not supported")
}
