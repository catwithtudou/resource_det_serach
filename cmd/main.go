package main

import (
	"flag"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"resource_det_search/internal/conf"
	"resource_det_search/internal/utils"
)

var (
	FlagConf string
)

func init() {
	flag.StringVar(&FlagConf, "conf", "../configs/config.yaml", "config path, eg: -conf config.yaml")
}

func InitLogger() *zap.SugaredLogger {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}

	writeSyncer := zapcore.AddSync(lumberJackLogger)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())

	return logger.Sugar()
}

func InitConf(flagConf string) *conf.Bootstrap {
	viper.SetConfigFile(flagConf)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("[InitConf]Fatal error config file: %s \n", err))
	}
	var cf *conf.Bootstrap
	if err := viper.Unmarshal(&cf); err != nil {
		panic(fmt.Errorf("[InitConf]Fatal error config structure: %s \n", err))
	}

	return cf
}

func main() {
	flag.Parse()
	logger := InitLogger()
	defer logger.Sync()
	cf := InitConf(FlagConf)

	utils.NewQny(cf.Data.Qiniuyun.Bucket, cf.Data.Qiniuyun.Access, cf.Data.Qiniuyun.Secret, cf.Data.Qiniuyun.Domain)

	app, _, err := initApp(cf.Data, logger)
	if err != nil {
		panic(fmt.Errorf("Fatal Init App: %s \n", err))
	}

	err = app.HttpEngine.Run(":" + cf.Server.Http.Port)
	if err != nil {
		panic(fmt.Errorf("Fatal Init HttpEngine: %s \n", err))
	}
}
