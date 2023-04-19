package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"os/signal"
	"stubs/internal"
	"syscall"
)

func main() {
	cfgFilePath := flag.String("config", "./config.yaml", "Configuration file path")
	flag.Parse()
	ctx, stop := context.WithCancel(context.Background())

	signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGILL, syscall.SIGINT)

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Application fatal error: %s", err)
			stop()
		}
	}()

	app, err := internal.NewApp(ctx, *cfgFilePath)
	if err != nil {
		logrus.Fatalf("Error on application starting: %s\n", err.Error())
	}
	if runErr := app.Run(); runErr != nil {
		if runErr != nil {
			return
		}
	}
}
