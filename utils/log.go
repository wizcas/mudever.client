package utils

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

// InitLogger with given configrations
func InitLogger(name string, production bool) {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = l.Named(name)
}

// Log returns a logger of the root namespace
func Logger() *zap.SugaredLogger {
	return logger.Sugar()
}
