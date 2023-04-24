package http

type Path struct {
	Name      string                  `yaml:"name"`
	Url       string                  `yaml:"url"`
	Delay     string                  `yaml:"delay"`
	Protocol  string                  `yaml:"protocol"`
	Headers   map[string]string       `yaml:"headers"`
	Method    string                  `yaml:"method"`
	Responses map[string]PathResponse `yaml:"responses"`
}

type PathResponse struct {
	When string           `yaml:"when"`
	Code int              `yaml:"code"`
	Body PathResponseBody `yaml:"body"`
}

type PathResponseBody struct {
	Template string            `yaml:"template"`
	Params   map[string]string `yaml:"params"`
	Store    map[string]string `yaml:"store"`
}
