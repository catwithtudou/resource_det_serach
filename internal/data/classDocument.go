package data

import (
	"context"
	"errors"
	"resource_det_search/internal/biz"
	"strconv"
)

type classDocumentRepo struct {
	data *Data
	idx  string
}

func NewClassDocumentRepo(data *Data) biz.IClassDocumentRepo {
	return &classDocumentRepo{
		data: data,
		idx:  "class_document",
	}
}

func (c *classDocumentRepo) InsertDoc(ctx context.Context, docId uint, cd *biz.ClassDocument) error {
	if docId <= 0 || cd == nil {
		return errors.New("docId or cd is nil")
	}

	_, err := c.data.es.Index().Index(c.idx).Id(strconv.Itoa(int(docId))).BodyJson(cd).Do(ctx)
	if err != nil {
		return err
	}

	_, err = c.data.es.Flush().Index(c.idx).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
