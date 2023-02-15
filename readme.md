[![Build Status](https://github.com/neggles/middleflare/workflows/TraefikTest/badge.svg?branch=master)](https://github.com/neggles/middleflare/actions)

# middleflare

a small traefik middleware plugin to remap CF-Connecting-IP over the top of X-Real-IP and X-Forwarded-For.

This has been done before but I didn't like the other implementations & wanted an excuse to do a golang thing.

Probably do not use this in production. Also, it looks like it maps CF-Visitor's "scheme" value, but it doesn't yet.

## Configuration

Define the plugin:

```yaml
# Static configuration

experimental:
  plugins:
    middleflare:
      moduleName: github.com/neggles/middleflare
      version: v0.0.2
```

### Example configuration:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - middleflare

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    middleflare:
      plugin:
        middleflare:
          includeDefault: true
          trustedProxies: []
```

### Kubernetes CRDs:

`TODO`: Add these when it's not 1:30am...
