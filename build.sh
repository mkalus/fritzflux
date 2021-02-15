#!/usr/bin/env bash

docker buildx create --use --driver docker-container
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t ronix/fritzflux . --push