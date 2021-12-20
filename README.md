# Collect some metrics from your zerotier
Just trying some stuff out for now. Mostly VL1 related.

## Latency to other nodes
<img width="850" alt="graph of peer latency" src="https://user-images.githubusercontent.com/11598/146835763-bb993e6d-57e3-44a5-9be7-f6ab5cd08b78.png">

## Count direct vs relayed peers
<img width="915" alt="graph of direct vs relayed peers" src="https://user-images.githubusercontent.com/11598/146835858-5291c360-2485-4bc1-8f9a-be6bccdb74d3.png">

## TODO
- [ ] make labels for PLANET, LEAF, MOON
- [ ] integrationy tests
- [ ] package it up, so I can install it more easily on my devices
  - [x] rudimentary linux and macos services in examples dir
- [ ] configurable if it logs individual node addresses or not, just planets for example.

## run
```
go run cmd/main.go
```

### args so far
- `--web.listen-address` default: ::19993
- `--web.telemetry-path` default: /metrics
- `--zerotier.authtoken-path` default: /var/lib/zerotier-one (linux), $HOME/Library/Application Support/ZeroTier/One/authtoken.secret (mac), TODO (windows)


