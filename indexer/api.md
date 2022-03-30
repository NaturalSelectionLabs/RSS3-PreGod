# API

## Base URL

On production:

```
http://pregod-indexer-api.pregod.traefik.mesh:3000
```

On dev:

```
http://localhost:8081 # depends on `indexer.server.http_port` in `config/config.*.json`
```

## Endpoints

### Get User Bio

```
GET /bio?proof=<proof>&platform_id=<platform_id>
```

For `platform_id`, see [RIP-1](https://rss3.io/protocol/RIPs/RIP-1.html#account-platform-list).

Examples:

- Twitter: `GET /bio?proof=diygod&platform_id=6`
- Misskey: `GET /bio?proof=song@misskey.io&platform_id=7`
- Jike: `GET /bio?proof=C05E4867-4251-4F11-9096-C1D720B41710&platform_id=8`

### Get User Profile

TODO
