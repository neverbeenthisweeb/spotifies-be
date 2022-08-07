package log

import (
	"context"
	"spotifies-be/web/constant"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	logger, err := zap.NewProduction(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.DPanicLevel), // Disable stack trace
	)
	if err != nil {
		panic(err)
	}

	Logger = logger
}

func Info(ctx context.Context, msg string) {
	var requestID string
	if id := ctx.Value(constant.CtxKeyRequestID); id != nil {
		requestID = id.(string)
	}
	Logger.Info(msg,
		zap.String("request_id", requestID),
	)
}

func Warn(ctx context.Context, msg string) {
	var requestID string
	if id := ctx.Value(constant.CtxKeyRequestID); id != nil {
		requestID = id.(string)
	}
	Logger.Warn(msg,
		zap.String("request_id", requestID),
	)
}

func Error(ctx context.Context, msg string, err error) {
	var requestID string
	if id := ctx.Value(constant.CtxKeyRequestID); id != nil {
		requestID = id.(string)
	}
	Logger.Error(msg,
		zap.String("request_id", requestID),
		zap.String("error", err.Error()),
	)
}
