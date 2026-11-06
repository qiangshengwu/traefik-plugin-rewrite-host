# traefik‑plugin‑rewrite‑host

Traefik custom middleware plugin to **override upstream `Host`**, mainly for S3 compatible services like SeaweedFS, MinIO.

> The built‑in Traefik `Headers` middleware cannot modify `req.Host` (the internal field used by Go http client when sending requests to backend).
> This plugin changes `req.Host` directly, which takes effect for upstream HTTP request.

## Features
- Override request `Host` for upstream backend
- Overwrite `X‑Forwarded‑Host` header
- Optional keep original host into `X‑Original‑Host`
- Input validation, reject full URL, only accept `host[:port]` format
- Safe: use `req.Clone()` avoid polluting original http.Request

## Limitations
- Only works as traefik **middleware plugin** (experimental plugins)
- Not for changing Host seen by client browser; only modify request sent to backend service.

## Configuration

### Static config (traefik.yml / traefik.toml)
Enable experimental plugin:
```yaml
experimental:
  plugins:
    rewrite-host:
      moduleName: github.com/qiangshengwu/traefik-plugin-rewrite-host
      version: v0.0.1
```

### Dynamic configuration example
```yaml
http:
  middlewares:
    rewrite-s3-host:
      plugin:
        rewrite-host:
          host: "devops-seaweedfs-s3:8333"

  routers:
    s3-router:
      rule: PathPrefix(`/buckets`)
      service: seaweedfs-s3
      middlewares:
        - rewrite-s3-host

  services:
    seaweedfs-s3:
      loadBalancer:
        servers:
          - url: "http://devops-seaweedfs-s3:8333"
```

## How it works

1. Clone the incoming `*http.Request` to prevent mutation of the original request.
2. Assign configured value to cloned request `.Host` (controls outgoing Host header to backend).
3. Set the `X‑Forwarded‑Host` header.
4. Pass the cloned request to the next middleware / upstream backend.

>
> To enable `X‑Original‑Host` that keeps the real client‑side host, use the extended `ServeHTTP` snippet in `plugin.go`.

## Local Development & Testing

Run unit tests:

```bash
go test -v ./...
```

## Important Notes

- Never modify the incoming `req` pointer in‑place without clone, Traefik reuses request objects internally.
- Required for S3‑like services which validate signature based on the `Host` value; modifying only header map will not work.

## License

MIT