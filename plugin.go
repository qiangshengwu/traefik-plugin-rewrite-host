package traefik_plugin_rewrite_host

import (
	"context"
	"net/http"
	"strings"
)

type Config struct {
	Host string `json:"host,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type S3Host struct {
	next http.Handler
	host string
	name string
}

func New(
	ctx context.Context,
	next http.Handler,
	config *Config,
	name string,
) (http.Handler, error) {
	if config == nil {
		return nil, http.ErrNotSupported
	}

	host := strings.TrimSpace(config.Host)
	if host == "" {
		return nil, http.ErrNotSupported
	}

	if strings.Contains(host, "://") || strings.ContainsAny(host, "/?#") {
		return nil, http.ErrNotSupported
	}

	return &S3Host{
		next: next,
		host: host,
		name: name,
	}, nil
}

func (s *S3Host) ServeHTTP(
	rw http.ResponseWriter,
	req *http.Request,
) {
	req.Host = s.host
	req.Header.Set("X-Forwarded-Host", s.host)
	s.next.ServeHTTP(rw, req)
}
