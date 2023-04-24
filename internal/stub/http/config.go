package http

type Config struct {
	Host        string `json:"host" yaml:"host"`
	Server      string `host:"server" yaml:"server"`
	Middlewares []string
	Log         struct {
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
