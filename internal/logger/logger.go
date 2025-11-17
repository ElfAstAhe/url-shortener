package logger

import (
	"github.com/ElfAstAhe/url-shortener/internal/config"
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string, stage string) error {
	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	conf := getConfigByStage(stage)
	conf.Level = zapLevel

	zl, err := conf.Build()
	if err != nil {
		return err
	}

	Log = zl

	return nil
}

func getConfigByStage(stage string) *zap.Config {
	var conf zap.Config
	switch stage {
	case config.ProjectStageProduction:
		conf = zap.NewProductionConfig()
	default:
		conf = zap.NewDevelopmentConfig()
	}

	return &conf
}
