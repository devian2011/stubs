package internal

import (
	"context"
	"github.com/valyala/fasthttp"
	"stubs/internal/server"
)

type App struct {
	ctx         context.Context
	cfg         *Configuration
	httpHandler *server.HttpServer
}

func NewApp(ctx context.Context, configFilePath string) (*App, error) {
	cfg, err := NewConfiguration(configFilePath)
	if err != nil {
		return nil, err
	}
	return &App{
		ctx:         ctx,
		cfg:         cfg,
		httpHandler: server.NewHttpServer(),
	}, nil
}

func (a *App) Run() error {
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
