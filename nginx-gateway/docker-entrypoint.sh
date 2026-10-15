#!/bin/sh
set -e

CONSUL_ADDR=${CONSUL_ADDR:-consul:8500}

# Render initial config
consul-template \
  -consul-addr="$CONSUL_ADDR" \
  -template="/etc/nginx/templates/nginx.conf.ctmpl:/etc/nginx/nginx.conf" \
  -once

# Start Nginx in background
nginx

# Run consul-template to watch for service changes and reload nginx
exec consul-template \
  -consul-addr="$CONSUL_ADDR" \
  -template="/etc/nginx/templates/nginx.conf.ctmpl:/etc/nginx/nginx.conf:nginx -s reload"
