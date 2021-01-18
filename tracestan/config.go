package tracestan

import (
	"strings"

	"github.com/nats-io/stan.go"
)

type Config struct {
	Addr      string
	ClientID  string
	ClusterID string
}

func (c *Config) addr() stan.Option {
	const prefix = "nats://"

	if !strings.HasPrefix(c.Addr, prefix) {
		c.Addr = prefix + c.Addr
	}

	return stan.NatsURL(c.Addr)
}
