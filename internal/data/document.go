package data

import (
	"context"
	"resource_det_search/internal/biz"
)

type documentRepo struct {
	data *Data
}

func NewDocumentRepo(data *Data) biz.IDocumentRepo {
	return &documentRepo{
		data: data,
	}
}

func (d *documentRepo) GetDocs(ctx context.Context) ([]*biz.Document, error) {
	return nil, nil
}
func (d *documentRepo) GetDocById(ctx context.Context, id uint) (*biz.Document, error) {
	return nil, nil
}
func (d *documentRepo) GetDocsByUid(ctx context.Context, uid uint) ([]*biz.Document, error) {
	return nil, nil
}
func (d *documentRepo) GetDocsWithDid(ctx context.Context, did uint) ([]*biz.Document, error) {
	return nil, nil
}
func (d *documentRepo) UpdateDocById(ctx context.Context, doc *biz.Document) error {
	return nil
}
func (d *documentRepo) DeleteDocById(ctx context.Context, doc *biz.Document) error {
	return nil
}
