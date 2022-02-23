package logger

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger/engine"
	"go.uber.org/zap"
)

/*
 * Use: If you want to use,
 * import this package in the corresponding go file and call the corresponding format
 * var logger = logger.Logger
 * or
 * var desugarredLogger = logger.DesugarredLogger
 */
var Logger *zap.SugaredLogger
var DesugarredLogger *zap.Logger

// Some simple encapsulations are made for the upper layer.
// The Sugared mode of the zap library is used by default.
// You can customize the encapsulation here to use other log libraries.
func Setup() error {
	var err error
	Logger, err = engine.InitZapLogger(config.Config.Logger)

	if err != nil {
		return err
	}

	DesugarredLogger = Logger.Desugar()

	return nil
}
