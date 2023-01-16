package http

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

type Config struct {
	Host string `json:"host" yaml:"host"`
	Log  struct {
		Provider string `yaml:"type" json:"type"`
		Path     string `json:"directory" yaml:"directory"`
	} `json:"log" yaml:"log"`
	Directories struct {
		Rules  string `json:"rules" yaml:"rules"`
		Static string `json:"static" yaml:"static"`
	} `json:"directories" yaml:"directories"`
	Ignore         []string          `json:"ignore" yaml:"ignore"`
	DefaultHeaders map[string]string `json:"defaultHeaders" yaml:"defaultHeaders"`
}

type Stub struct {
	cfg   *Config
	rules *RuleSet
}

func NewHttpStub(cfg *Config) *Stub {
	stub := &Stub{
		cfg:   cfg,
		rules: newRuleSet(),
	}
	stub.rules.ScanDir(cfg.Directories.Rules)

	return stub
}

func (s *Stub) GetDomain() string {
	return s.cfg.Host
}

func (s *Stub) Process(ctx *fasthttp.RequestCtx) {
	rule, err := s.rules.FindPath(ctx.Path())
	if err != nil {
		ctx.Error(fmt.Sprintf("Unknown path: %s for host: %s", ctx.Path(), ctx.Host()), 404)
		return
	}
	time.Sleep(rule.Delay)

	for defaultHKey, defaultHValue := range s.cfg.DefaultHeaders {
		ctx.Response.Header.Add(defaultHKey, defaultHValue)
	}
	for hKey, hValue := range rule.Headers {
		ctx.Response.Header.Add(hKey, hValue)
	}
	httpCode, _ := strconv.Atoi(rule.HttpCode)
	ctx.Response.SetStatusCode(httpCode)
	ctx.WriteString(rule.Body)
}
