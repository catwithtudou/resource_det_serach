package data

import (
	"go.uber.org/zap"
	"resource_det_search/internal/conf"
	"testing"
)

func newData(t *testing.T) (*Data, *zap.SugaredLogger) {
	logger, _ := zap.NewProduction()
	reLogger := logger.Sugar()

	data, _, err := NewData(&conf.Data{
		Database: &conf.Data_Database{
			Driver: "mysql",
			Source: "root:a949812478@tcp(127.0.0.1:3306)/resource_det_search?charset=utf8&parseTime=True&loc=Local",
		},
		Es: &conf.Data_Es{
			Source: "http://localhost:9200",
		},
	}, reLogger)
	if err != nil {
		t.Fatal(err)
	}

	return data, reLogger
}

func TestNewData(t *testing.T) {
	_, _ = newData(t)
}
