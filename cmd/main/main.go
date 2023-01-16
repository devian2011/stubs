package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"stubs/internal"
	"syscall"
)

func main() {
	cfgFilePath := flag.String("config", "./config.yaml", "Configuration file path")
	flag.Parse()
	ctx, stop := context.WithCancel(context.Background())

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGILL, syscall.SIGINT)
	go func() {
		<-signalCh
		stop()
	}()

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
