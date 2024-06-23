package lib

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/goProjects/loan_app/lib/logger"
	"github.com/goProjects/loan_app/lib/utils"
	"go.uber.org/zap"
)

// Executing a call to recover inside a deferred function stops the panicking sequence,
// hence adding the recover block here. Use errHandler to handle the error
// created due to panic recovery, as per use case.
func PanicHandler(ctx context.Context, logFields map[string]any, errHandler func(error)) {
	if r := recover(); r != nil {
		var e error
		if m, ok := r.(string); ok {
			e = errors.New(m)
		} else if m, ok := r.(error); ok {
			e = m
		} else {
			e = errors.New(utils.ConvertToString(r))
		}

		if logFields == nil {
			logFields = map[string]any{}
		}

		logFields["error"] = e.Error()
		logFields["stack_trace"] = string(debug.Stack())
		logger.E(ctx, "recovered from panic", zap.Any("payload", logFields))

		if errHandler != nil {
			errHandler(e)
		}
	}
}
