package config

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	MaxIdleConns        int           `split_words:"true" default:"10"`
	MaxIdleConnsPerHost int           `split_words:"true" default:"2"`
	IdleConnTimeout     time.Duration `split_words:"true" default:"30s"`
	Timeout             time.Duration `split_words:"true" default:"10s"`
}

func (h HTTPClient) NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        h.MaxIdleConns,
			MaxIdleConnsPerHost: h.MaxIdleConnsPerHost,
			IdleConnTimeout:     h.IdleConnTimeout,
		},
		Timeout: h.Timeout,
	}
}
