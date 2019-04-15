package zaplog

import (
	"go.uber.org/zap"
)

func InitLog() {

	parentLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(parentLogger)

}

// Sync call all logger's method `Sync`
func Sync() {
	zap.L().Sync()
}
