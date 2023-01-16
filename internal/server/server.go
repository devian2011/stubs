package server

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"strings"
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

	domain := strings.Split(string(ctx.Host()), ":")[0]
	if p, exists := h.processors[domain]; exists {
		path := ctx.Path()
		if string(path) == "/" || string(path) == "" {
			ctx.URI().SetPath("/index")
		}

		p.Process(ctx)
	} else {
		ctx.Error(fmt.Sprintf("Processor for domain: %s does not exists", domain), 404)
	}
}
