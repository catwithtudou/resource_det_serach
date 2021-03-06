package data

import (
	"errors"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/conf"
	"time"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var ProvideSet = wire.NewSet(NewData, NewUserRepo, NewDimensionRepo, NewDocumentRepo, NewClassDocumentRepo)

// Data .
type Data struct {
	db *gorm.DB
	es *elastic.Client
}

type GormWriter struct {
	logger *zap.SugaredLogger
}

func (g GormWriter) Printf(format string, args ...interface{}) {
	g.logger.Warnf(format, args)
}

func NewData(conf *conf.Data, log *zap.SugaredLogger) (*Data, func(), error) {
	// 1. db
	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		// 全局模式：执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		PrepareStmt: true,
		Logger: logger.New(
			GormWriter{logger: log},
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	})

	if err != nil {
		log.Errorf("[NewData]failed to open mysql resource:err=[%+v]", err)
		return nil, nil, err
	}

	if db == nil {
		err = errors.New("[NewData]db engine is nil")
		log.Errorf(err.Error())
		return nil, nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Errorf("[NewData]failed to open connect pool:err=[%+v]", err)
		return nil, nil, err
	}

	sqlDb.SetMaxIdleConns(25)
	sqlDb.SetMaxOpenConns(25)
	sqlDb.SetConnMaxLifetime(5 * time.Minute)

	// create table
	_ = db.AutoMigrate(&biz.User{}, &biz.Dimension{}, &biz.Document{}, &biz.DocWithDm{})

	// 2. es
	es, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(conf.Es.Source))
	if err != nil {
		log.Errorf("[NewData]failed to new es client:err=[%+v]", err)
		return nil, nil, err
	}

	return &Data{db: db, es: es}, func() {}, nil
}
