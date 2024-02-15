package internal

import (
	"market/global"
	"os"

	"go.uber.org/zap/zapcore"
)

type fileRotatelogs struct{}

var FileRotatelogs = new(fileRotatelogs)

func (r *fileRotatelogs) GetWriteSyncer(level string) zapcore.WriteSyncer {
	fileWriter := NewCutter(global.MARKET_CONFIG.Zap.Director, level, WithCutterFormat("2006-01-02"))
	if global.MARKET_CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	}
	return zapcore.AddSync(fileWriter)
}
