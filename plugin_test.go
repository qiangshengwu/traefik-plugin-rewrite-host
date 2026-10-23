package traefik_plugin_rewrite_host

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRewriteHost(t *testing.T) {
	var gotHost string

	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		gotHost = r.Host
		w.WriteHeader(http.StatusNoContent)
	})

	handler, err := New(context.Background(), next, &Config{
		Host: "devops-seaweedfs-s3:8333",
	}, "test")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"http://192.168.31.136/buckets/test.png",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if gotHost != "devops-seaweedfs-s3:8333" {
		t.Fatalf("unexpected Host: %q", gotHost)
	}
}
