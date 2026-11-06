# traefik‑plugin‑rewrite‑host

Traefik 自定义中间件插件，用于重写发往上流后端的 HTTP `Host`，主要适配 SeaweedFS、MinIO 这类 S3 兼容后端服务。

> Traefik 官方内置的 `Headers` 中间件仅能操作请求头 Map，**无法修改 Go HTTP 内部 `req.Host` 字段**，因此对 S3 类服务不生效。
> 本插件在克隆后的请求对象上直接修改 `req.Host`，真正改变发往后端的 Host 请求头。

## 功能特性
- 重写发往上流后端的请求 `Host`
- 重写 `X‑Forwarded‑Host` 请求头
- 可选：将原始客户端 Host 保存至 `X‑Original‑Host`（参考源码注释）
- 严格参数校验：拒绝完整URL，仅接受 `host[:port]` 格式
- 安全实现：使用 `req.Clone()`，不会污染 Traefik 内部复用的原始 `*http.Request`

## 使用限制
- 属于**私有实验性 Traefik 插件**，未上架 Traefik 公共插件市场。
- 仅修改发往后端的请求，**不会改变终端浏览器看到的访问域名**。

## 配置说明

### 静态配置 traefik.yml
开启实验插件能力：

```yaml
experimental:
  plugins:
    rewrite-host:
      moduleName: github.com/qiangshengwu/traefik-plugin-rewrite-host
      version: v0.0.1
```

### 动态配置示例
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

## 参数列表

| 参数 | 必填 | 说明 |
| ---- | ---- | ---- |
| `host` | 是 | 上游目标主机，格式 `hostname` 或 `hostname:port`。**禁止填写完整URL，不要带 http://、路径、查询参数。** |

## 本地开发 & 单元测试
执行单元测试
```bash
go test -v ./...
```

## 重要提醒
- Traefik 内部会复用请求对象，**禁止不做 Clone 直接原地修改入参 req**。
- S3 类服务会基于 Host 做签名校验，只修改 Header Map 不会生效，必须修改内部 req.Host。

## License

MIT