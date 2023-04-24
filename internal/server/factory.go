package server

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrHttpServerWithSameNameAlreadyExists = errors.New("server with same name already exists")
	ErrUnknownHttpServer                   = errors.New("unknown http server")
	ErrServerMustBeImplementRunnable       = errors.New("server must be implement runnable")
)

type Runnable interface {
	Run(context.Context, chan error)
}

type Store struct {
	ctx         context.Context
	errCh       chan error
	httpServers map[string]Runnable
}

func InitStore(ctx context.Context, cfgs []Config) (*Store, error) {
	factory := &Store{
		ctx:         ctx,
		errCh:       make(chan error, len(cfgs)),
		httpServers: make(map[string]Runnable, 0),
	}
	for _, cfg := range cfgs {
		switch cfg.Type {
		case "http", "jsonrpc":
			if _, exists := factory.httpServers[cfg.Name]; exists {
				return nil, ErrHttpServerWithSameNameAlreadyExists
			}
			httpServer := NewHttpServer(cfg.Addr, cfg.Params)
			if _, ok := interface{}(*httpServer).(Runnable); !ok {
				return nil, ErrServerMustBeImplementRunnable
			}
		}
	}

	return factory, nil
}

func (s *Store) ServersRun() {
	go func() {
		for _, server := range s.httpServers {
			server.Run(s.ctx, s.errCh)
		}
	}()
	go func() {
		for e := range s.errCh {
			logrus.Errorf("error on server: %s", e.Error())
		}
	}()
}

func (s *Store) GetHttp(name string) (*HttpServer, error) {
	if http, exists := s.httpServers[name]; exists {
		return http.(*HttpServer), nil
	}

	return nil, ErrUnknownHttpServer
}
