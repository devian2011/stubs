package stub

import "stubs/internal/stub/http"

type Config struct {
	Http http.Config `json:"http" yaml:"http"`
}
