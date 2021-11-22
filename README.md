# Collect some metrics from your zerotier
Just trying some stuff out for now

## TODO
- [] integration test
- [] package it up, so I can install it more easily on my devices
- [] configurable if it logs individual node addresses or not

## run
```
go run cmd/main.go
```

### args so far
- `--web.listen-address`
- `--web.telemetry-path`
- `--zerotier.authtoken-path`
