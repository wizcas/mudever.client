package log

import (
	"github.com/wizcas/mudever.svc/utils"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// Logger returns a logger of the 'terminal' namespace
func Logger() *zap.SugaredLogger {
	if logger == nil {
		logger = utils.Logger().Named("terminal")
	}
	return logger
}
