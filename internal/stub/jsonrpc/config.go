package http

type Config struct {
	Host        string `json:"host" yaml:"host"`
	Server      string `json:"server" yaml:"server"`
	Endpoint    string `json:"endpoint" yaml:"endpoint"`
	Middlewares []string
	Log         struct {
		Provider string `yaml:"type" json:"type"`
		Path     string `json:"directory" yaml:"directory"`
	} `json:"log" yaml:"log"`
	Rules string `json:"rules" yaml:"rules"`
}
