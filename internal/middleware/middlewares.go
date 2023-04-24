package middleware

type Config struct {
	Name   string
	Plugin string
	Params []string
}

type Store struct {
}

func NewStore(cfg []Config) (*Store, error) {
	return &Store{}, nil
}

