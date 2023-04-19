package server

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net"
)

type Processor interface {
	GetDomain() string
	Process(ctx *fasthttp.RequestCtx)
}

type Middleware interface {
	ForRequest() bool
	ForResponse() bool
	Handle(ctx *fasthttp.RequestCtx)
}

type Handler struct {
	middleware []Middleware
	processors map[string]Processor
}

func NewHandler() *Handler {
	return &Handler{
		processors: make(map[string]Processor, 0),
		middleware: make([]Middleware, 0),
	}
}

func (h *Handler) AddMiddleware(m Middleware) {
	h.middleware = append(h.middleware, m)
}

func (h *Handler) AddProcessor(processor Processor) error {
	if _, exists := h.processors[processor.GetDomain()]; exists {
		return errors.New(fmt.Sprintf("processor for domain: %s already exists", processor.GetDomain()))
	}
	h.processors[processor.GetDomain()] = processor
	return nil
}

func (h *Handler) Handle(ctx *fasthttp.RequestCtx) {
	host, _, err := net.SplitHostPort(string(ctx.Host()))
	if err != nil {
		logrus.Errorf("cannot find host and port. host: %s, err: %s", ctx.Host(), err.Error())
		ctx.Error(fmt.Sprintf("cannot find host or port in host:%s", host), 500)
	}

	if p, exists := h.processors[host]; exists {
		path := ctx.Path()
		if string(path) == "/" || len(path) == 0 {
			ctx.URI().SetPath("/index")
		}

		p.Process(ctx)
	} else {
		ctx.Error(fmt.Sprintf("Processor for domain: %s does not exists", host), 404)
	}
}
