#!/bin/bash -e

TRAEFIK_HTTP_PORT="8080"
TRAEFIK_DASHBOARD_HTTP_PORT="8000"

mkdir -p app/bin
docker run --rm -v "${PWD}/app":/app -e CGO_ENABLED=0 -e GOOS=linux -w /app golang:1.6 go build -a -x -o bin/html5sse

if [[ ! -f "app/bin/traefik" ]]; then
  curl -L https://github.com/containous/traefik/releases/download/v1.0.0/traefik_linux-amd64 -o app/bin/traefik
  chmod 755 app/bin/traefik
fi

if [[ ! -f "app/bin/forego" ]]; then
  cd app/bin
  curl -L https://bin.equinox.io/c/ekMN3bCZFUn/forego-stable-linux-amd64.tgz | tar xzf -
  chmod 755 forego
  cd -
fi

docker build -t ashmckenzie/traefik_html5sse_example .

docker run --rm -ti \
-p "${TRAEFIK_HTTP_PORT}:${TRAEFIK_HTTP_PORT}" \
-p "${TRAEFIK_DASHBOARD_HTTP_PORT}:${TRAEFIK_DASHBOARD_HTTP_PORT}" \
-e "TRAEFIK_HTTP_PORT=${TRAEFIK_HTTP_PORT}" \
-e "TRAEFIK_DASHBOARD_HTTP_PORT=${TRAEFIK_DASHBOARD_HTTP_PORT}" \
ashmckenzie/traefik_html5sse_example
