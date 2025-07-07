package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Log = logger
}
