package options

import "time"

// ref: https://youthlin.com/20201762.html

type Client struct {
	Timeout time.Duration
	Cluster string
}

type Config struct {
	Timeout time.Duration
	Cluster string
}

type Option func(*Config)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithCluster(cluster string) Option {
	return func(c *Config) {
		c.Cluster = cluster
	}
}

func New(opts ...Option) *Client {
	config := &Config{
		Timeout: 20 * time.Millisecond,
		Cluster: "default",
	}

	for _, opt := range opts {
		opt(config)
	}

	return &Client{
		Timeout: config.Timeout,
		Cluster: config.Cluster,
	}
}
