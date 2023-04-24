package server

import (
	"context"
	"github.com/valyala/fasthttp"
)

type HttpServer struct {
	fHttp  *fasthttp.Server
	addr   string
	params []string
}

func NewHttpServer(addr string, params []string) *HttpServer {
	server := &HttpServer{
		fHttp: &fasthttp.Server{
			Handler: handle,
			Name:    addr,
		},
		addr:   "",
		params: params,
	}

	return server
}

func (h *HttpServer) Run(ctx context.Context, errCh chan error) {
	go func() {
		errStart := h.fHttp.ListenAndServe(h.addr)
		if errStart != nil {
			errCh <- errStart
			return
		}
		<-ctx.Done()
		errStop := h.fHttp.Shutdown()
		if errStop != nil {
			errCh <- errStop
			return
		}
	}()
}

func handle(ctx *fasthttp.RequestCtx) {

}
