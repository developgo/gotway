# microgateway
[![Build Status](https://travis-ci.org/gosmo-devs/microgateway.svg)](https://travis-ci.org/gosmo-devs/microgateway)

A simple, lightweight and blazingly fast API gateway written in Go

## Development

```
$ ./scripts/run-dev.sh
```

Run services in your machine using [tmux](https://github.com/tmux/tmux/wiki). Additional services, like Redis, are defined in [docker-compose.dev.yml](./docker-compose.dev.yml).

## Build Docker images locally

```
$ ./scripts/build.sh
```

## Run services using local images

```
$ docker-compose up -d
```
This can be useful for testing images before pushing them to DockerHub.

## Release images for multiple architectures

```
$ ./scripts/release.sh
```
This script will be executed by TravisCI when a tag is pushed.

## Deploy to a Kubernetes cluster

```
$ ./scripts/deploy.sh <environment>
```

## Services

|Service|DockerHub|PRE|
|-------|---------|---|
|Microgateway|[gosmogolang/microgateway](https://hub.docker.com/r/gosmogolang/microgateway)|https://pre-microgateway.gosmo-devs.duckdns.org|
|Catalog|[gosmogolang/catalog](https://hub.docker.com/r/gosmogolang/catalog)|https://pre-catalog.gosmo-devs.duckdns.org|
|Stock|[gosmogolang/stock](https://hub.docker.com/r/gosmogolang/stock)|https://pre-stock.gosmo-devs.duckdns.org|
