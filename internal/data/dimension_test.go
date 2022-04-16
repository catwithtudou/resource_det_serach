package data

import (
	"context"
	"resource_det_search/internal/utils"
	"testing"
)

func newDimensionRepoTest(t *testing.T) (*dimensionRepo, context.Context) {
	data, _ := newData(t)
	return &dimensionRepo{data: data}, context.Background()
}

func TestGetDmByDidUid(t *testing.T) {
	d, ctx := newDimensionRepoTest(t)
	result, err := d.GetDmByDidUid(ctx, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(utils.JsonToString(result))
}

func TestGetDmsByType(t *testing.T) {
	//d, ctx := newDimensionRepoTest(t)
	//result, err := d.GetDmsByType(ctx, 3, "tag")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//for _, v := range result {
	//	t.Logf(utils.JsonToString(v))
	//}
}
