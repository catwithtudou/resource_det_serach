package usecase

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mime/multipart"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/constants"
	"resource_det_search/internal/utils"
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
	err = d.checkDmIdsIllegal(ctx, doc.Uid, constants.Category, categories)
	if err != nil {
		return constants.DefaultErr, fmt.Errorf("[UploadUserDocument]illegal category id:err=[%+v],ids=[%+v]", err, categories)
	}
	err = d.checkDmIdsIllegal(ctx, doc.Uid, constants.Tag, tags)
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

	// todo:async goroutine handle:det file details (with update the database) and upload the es search

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

	// 直接识别部分
	if utils.DetByteTypesContains(fileType) {
		detType := constants.DetByteType(fileType)
		switch detType {
		case constants.Txt:
			return string(fileBytes), nil
		case constants.Docx:
			txt, err := utils.DetDocxByUnidoc(fileBytes)
			if err != nil {
				return "", fmt.Errorf("[DetFile]failed to DetDocxByUnidoc:err=[%+v]", err)
			}
			return txt, nil
		case constants.Pptx:
			txt, err := utils.DetPptxByUnidoc(fileBytes)
			if err != nil {
				return "", fmt.Errorf("[DetFile]failed to DetPptxByUnidoc:err=[%+v]", err)
			}
			return txt, nil
		case constants.Xlsx:
			txt, err := utils.DetXlsxByUnidoc(fileBytes)
			if err != nil {
				return "", fmt.Errorf("[DetFile]failed to DetXlsxByUnidoc:err=[%+v]", err)
			}
			return txt, nil
		case constants.Md:
			return utils.DetMd(fileBytes)
		}

	}

	// OCR识别部分
	if utils.DetOcrTypesContains(fileType) {

	}

	return "", nil
}

func (d *documentUsecase) uploadSearch(ctx context.Context, doc *biz.Document) error {
	fileType := doc.Type

	// 直接识别部分
	if fileType == "" {

	}

	// OCR识别部分
	if fileType == "" {

	}

	return nil
}

func (d *documentUsecase) checkDmIdsIllegal(ctx context.Context, uid uint, typeStr constants.DmType, ids []uint) error {
	if len(ids) <= 0 {
		return nil
	}

	result, err := d.dmRepo.GetUidTypeInIds(ctx, ids)
	if err != nil {
		return err
	}
	if len(result) != 1 || result[0].Uid != uid || result[0].Type != string(typeStr) {
		return errors.New("illegal dm id")
	}

	return nil
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
