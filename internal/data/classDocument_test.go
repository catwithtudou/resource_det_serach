package data

import (
	"context"
	"resource_det_search/internal/constants"
	"testing"
)

func newClassDocumentRepoTest(t *testing.T) (*classDocumentRepo, context.Context) {
	data, _ := newData(t)
	return &classDocumentRepo{
		data: data,
		idx:  constants.ClassDocument,
	}, context.Background()
}

func TestUpdateNums(t *testing.T) {
	c, ctx := newClassDocumentRepoTest(t)
	err := c.UpdateNums(ctx, 18, 1, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
}
