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

### Get Items

```
GET /item?proof=<proof>&platform_id=<platform_id>&network_id=<network_id>&limit=<limit>&timestamp=<timestamp>
```

- `proof`: (required) the proof of the user.
- `platform_id`: (required) see [RIP-1](https://rss3.io/protocol/RIPs/RIP-1.html#account-platform-list).
- `network_id`: (required) see [RIP-3](https://rss3.io/protocol/RIPs/RIP-3.html#item-network-list).
- `owner_id`: (required) the proof of the main account.
- `owner_platform_id`: (required) see [RIP-1](https://rss3.io/protocol/RIPs/RIP-1.html#account-platform-list).
- `limit`: (optional) the number of items to return. default: 100.
- `timestamp`: (optional) the timestamp of the last item. default: now.

Examples:

- Ethereum platforms:
  - Ethereum: `GET /item?proof=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&platform_id=1&network_id=1&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
  - Polygon: `GET /item?proof=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&platform_id=1&network_id=2&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
  - Bnb: `GET /item?proof=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&platform_id=1&network_id=3&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
  - Arbitrum: `GET /item?proof=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&platform_id=1&network_id=4&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
  - POAP: `GET /item?proof=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&platform_id=1&network_id=7&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
  - Mirror (Arweave): `GET /item?proof=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&platform_id=1&network_id=10&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
  - Gitcoin: ``
- Twitter: `GET /item?proof=diygod&platform_id=6&network_id=12&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
- Misskey: `GET /item?proof=Candinya@nya.one&platform_id=7&network_id=13&limit=10&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
- Jike: `GET /item?proof=169a5be6-f874-4df9-a2d1-a07e9e0e429b&platform_id=8&network_id=14&owner_id=0x08d66b34054a174841e2361bd4746Ff9F4905cC2&owner_platform_id=1&profile_source_id=0`
