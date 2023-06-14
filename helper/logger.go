package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// Topic for setting topic of log
	Topic = "book-store-log"
)

var (
	service string
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func LogContext(ctx context.Context, c, s string) *log.Entry {
	return log.WithFields(log.Fields{
		"topic":   Topic,
		"context": c,
		"scope":   s,
		"service": service,
	})
}

func Log(ctx context.Context, level log.Level, err error, context, scope string) {
	log.SetFormatter(&log.JSONFormatter{})
	/*
		syslogOutput, err := logrusSyslog.NewSyslogHook("", "", syslog.LOG_INFO, LogTag)
		// avoiding nil pointer. check for error first then add hook.
		if err != nil {
			log.Debug("Unable to setup syslog output")
		}
		log.AddHook(syslogOutput)
	*/
	entry := LogContext(ctx, context, scope)
	reqID := ctx.Value(echo.HeaderXRequestID)
	if reqID != "" {
		entry = entry.WithField("request_id", reqID)
	}
	message := err.Error()

	var sb strings.Builder
	sb.WriteString(message)
	sb.WriteString("\n")
	// i:=0
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			sb.WriteString(fmt.Sprintf("%+s:%d\n", f, f))
			// i++
			// if i == 3{
			// 	break
			// }
		}
	}

	message = sb.String()

	switch level {
	case log.DebugLevel:
		entry.Debug(message)
	case log.InfoLevel:
		entry.Info(message)
	case log.WarnLevel:
		entry.Warn(message)
	case log.ErrorLevel:
		entry.Error(message)
	case log.FatalLevel:
		entry.Fatal(message)
	case log.PanicLevel:
		entry.Panic(message)
	}
}
