package common

import (
	"go.uber.org/zap"
)

// Logger returns a logger of the 'terminal' namespace
func Logger() *zap.SugaredLogger {
	return zap.S().Named("nvt")
}
