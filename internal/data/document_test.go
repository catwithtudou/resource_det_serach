package data

import (
	"context"
	"gorm.io/gorm"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/utils"
	"testing"
)

func newDocumentRepoTest(t *testing.T) (*documentRepo, context.Context) {
	data, _ := newData(t)
	return &documentRepo{data: data}, context.Background()
}

func TestGetDocs(t *testing.T) {
	d, ctx := newDocumentRepoTest(t)
	result, err := d.GetDocs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range result {
		t.Logf(utils.JsonToString(v))
	}
}
func TestGetDocById(t *testing.T) {
	d, ctx := newDocumentRepoTest(t)
	result, err := d.GetDocById(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(utils.JsonToString(result))
}
func TestGetDocsByUid(t *testing.T) {
	d, ctx := newDocumentRepoTest(t)
	result, err := d.GetDocsByUid(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range result {
		t.Logf(utils.JsonToString(v))
	}
}
func TestGetDocsWithDid(t *testing.T) {
	d, ctx := newDocumentRepoTest(t)
	result, err := d.GetDocsWithDid(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range result {
		t.Logf(utils.JsonToString(v))
	}
}
func TestUpdateDocById(t *testing.T) {
	d, ctx := newDocumentRepoTest(t)
	err := d.UpdateDocById(ctx, &biz.Document{
		Model:        gorm.Model{ID: 1},
		Intro:        "about math",
		Title:        "math_docx",
		DownloadNum:  1,
		ScanNum:      1,
		LikeNum:      1,
		IsLoadSearch: true,
		IsSave:       true,
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestAddDocLikeNum(t *testing.T) {
	d, ctx := newDocumentRepoTest(t)
	err := d.AddDocLikeNum(ctx, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDocById(t *testing.T) {
	d, ctx := newDocumentRepoTest(t)
	err := d.DeleteDocById(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDocByIdWithUid(t *testing.T) {
	d, ctx := newDocumentRepoTest(t)
	err := d.DeleteDocByIdWithUid(ctx, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
}
