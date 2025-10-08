package config

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	MaxIdleConns        int           `default:"10"`
	MaxIdleConnsPerHost int           `default:"2"`
	IdleConnTimeout     time.Duration `default:"30s"`
	Timeout             time.Duration `default:"10s"`
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
