package consulutil

import "github.com/hashicorp/consul/api"

type Config struct {
	Addr   string `envconfig:"CONSUL_ADDR"`
	Scheme string `envconfig:"CONSUL_SCHEME"`
}

func Client(c Config) (*api.Client, error) {
	cfg := api.DefaultConfig()
	if c.Addr != "" {
		cfg.Address = c.Addr
	}
	if c.Scheme != "" {
		cfg.Scheme = c.Scheme
	}
	return api.NewClient(cfg)
}
