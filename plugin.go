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

	// Host 必须是主机名或主机名:端口，
	// 不允许配置 URL、路径或查询参数。
	if strings.Contains(host, "://") ||
		strings.ContainsAny(host, "/?#") {
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
	// 修改请求的实际 Host 字段，
	// 而不是 req.Header["Host"]。
	req.Host = s.host

	s.next.ServeHTTP(rw, req)
}
