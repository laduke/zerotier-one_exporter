# Collect some metrics from your zerotier
Just trying some stuff out for now

## TODO
- [] integrationy tests
- [] package it up, so I can install it more easily on my devices
- [] configurable if it logs individual node addresses or not

## run
```
go run cmd/main.go
```

### args so far
- `--web.listen-address` default: ::19993
- `--web.telemetry-path` default: /metrics
- `--zerotier.authtoken-path` default: /var/lib/zerotier-one (linux), $HOME/Library/Application Support/ZeroTier/One/authtoken.secret (mac), TODO (windows)

