package main

import (
	"context"

	"github.com/amarin/logging"
)

type App struct {
	logging.Logger
}

// doSomethingCtx demonstrates how to use context loggers.
func doSomethingCtx(ctx context.Context) {
	// use default context-aware logger
	log := logging.NewLoggerCtx(ctx)
	log.Infof("From inside doSomethingCtx")

	// use named context-aware logger
	namedLogger := logging.NewNamedLoggerCtx(ctx, "named")
	namedLogger.Infof("Named from inside doSomethingCtx")
}

func main() {
	// init logging subsystem
	logging.MustInit()

	// use logging as explicit logger instance
	logger := logging.NewLogger()
	logger.Info("simple new logger")

	// use embedded logger
	app := &App{Logger: logging.NewNamedLogger("app", logging.LevelInfo)}
	app.Info("embedded app logger")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doSomethingCtx(ctx)
}
