package internal

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"stubs/internal/server"
	"stubs/internal/stub/http"
)

type App struct {
	ctx         context.Context
	cfg         *Configuration
	httpHandler *server.Handler
}

func NewApp(ctx context.Context, configFilePath string) (*App, error) {
	cfg, err := NewConfiguration(configFilePath)
	if err != nil {
		return nil, err
	}
	return &App{
		ctx:         ctx,
		cfg:         cfg,
		httpHandler: server.NewHandler(),
	}, nil
}

func (a *App) loadHttpStubs() {
	for _, stubConf := range a.cfg.Stubs.Http {
		stub := http.NewHttpStub(&stubConf)
		addStubProcErr := a.httpHandler.AddProcessor(stub)
		if addStubProcErr != nil {
			logrus.Errorf("Error on add stub processor for: %s err: %s", stub.GetDomain(), addStubProcErr.Error())
		}
	}
}

func (a *App) Run() error {
	a.loadHttpStubs()

	srv := fasthttp.Server{
		Handler: a.httpHandler.Handle,
	}

	errCh := make(chan error)

	go func(addr string, errCh chan error) {
		err := srv.ListenAndServe(a.cfg.Http.Addr)
		if err != nil {
			errCh <- err
		}
	}(a.cfg.Http.Addr, errCh)

	select {
	case <-a.ctx.Done():
		srv.Shutdown()
		return nil
	case err := <-errCh:
		return err
	}
}
