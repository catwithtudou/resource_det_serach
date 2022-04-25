package data

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"reflect"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/constants"
	"strconv"
)

type classDocumentRepo struct {
	data *Data
	idx  string
}

func NewClassDocumentRepo(data *Data) biz.IClassDocumentRepo {
	return &classDocumentRepo{
		data: data,
		idx:  constants.ClassDocument,
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

func (c *classDocumentRepo) SearchAllQuery(ctx context.Context, queryStr string, sortBy string) ([]*biz.ClassDocument, error) {
	if queryStr == "" {
		return nil, errors.New("query str is nil")
	}

	query := elastic.NewQueryStringQuery(queryStr)
	res, err := c.data.es.Search(c.idx).Query(query).Sort(sortBy, false).Do(ctx)
	if err != nil {
		return nil, err
	}

	return c.searchCDValue(res), nil
}

func (c *classDocumentRepo) SearchQueryByPart(ctx context.Context, queryStr string, partName string, sortBy string) ([]*biz.ClassDocument, error) {
	if queryStr == "" || partName == "" {
		return nil, errors.New("queryStr or partName is nil")
	}

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewMatchPhraseQuery("part", partName), elastic.NewQueryStringQuery(queryStr))
	res, err := c.data.es.Search(c.idx).Query(query).Sort(sortBy, false).Do(ctx)
	if err != nil {
		return nil, err
	}

	return c.searchCDValue(res), nil
}

func (c *classDocumentRepo) UpdateNums(ctx context.Context, docId uint, likeNum uint, scanNum uint, downloadNum uint) error {
	if docId <= 0 {
		return errors.New("docId is nil")
	}

	_, err := c.data.es.Update().Index(c.idx).Id(strconv.Itoa(int(docId))).Doc(map[string]interface{}{
		"like_num": likeNum, "scan_num": scanNum, "download_num": downloadNum,
	}).Refresh("true").Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *classDocumentRepo) UpdateDimensions(ctx context.Context, docId uint, typeStr constants.DmType, oldDmName string, newDmName string) error {
	if docId <= 0 || typeStr == "" || oldDmName == "" || newDmName == "" {
		return errors.New("docId or typeStr or oldDmName or newDmName is nil")
	}

	res, err := c.data.es.Get().Index(c.idx).Id(strconv.Itoa(int(docId))).Do(ctx)
	if err != nil {
		return err
	}

	var cd biz.ClassDocument
	err = json.Unmarshal(res.Source, &cd)
	if err != nil {
		return err
	}

	switch typeStr {
	case constants.Tag:
		for k, v := range cd.Tags {
			if v == oldDmName {
				cd.Tags[k] = newDmName
				break
			}
		}
		_, err = c.data.es.Update().Index(c.idx).Id(strconv.Itoa(int(docId))).Doc(map[string]interface{}{
			"tags": cd.Tags,
		}).Refresh("true").Do(ctx)
		if err != nil {
			return err
		}

		return nil
	case constants.Category:
		for k, v := range cd.Categories {
			if v == oldDmName {
				cd.Categories[k] = newDmName
				break
			}
		}
		_, err = c.data.es.Update().Index(c.idx).Id(strconv.Itoa(int(docId))).Doc(map[string]interface{}{
			"categories": cd.Categories,
		}).Refresh("true").Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (c *classDocumentRepo) searchCDValue(res *elastic.SearchResult) []*biz.ClassDocument {
	if res == nil {
		return nil
	}
	docs := make([]*biz.ClassDocument, 0)
	for _, doc := range res.Each(reflect.TypeOf(biz.ClassDocument{})) {
		if res, ok := doc.(biz.ClassDocument); ok {
			docs = append(docs, &res)
		}
	}
	return docs
}
