#!/bin/bash

export REGISTRY=ghcr.io
export ORG=awakari
export COMPONENT=core-tests
export SLUG=${REGISTRY}/${ORG}/${COMPONENT}
export VERSION=latest
docker tag awakari/core-tests "${SLUG}":"${VERSION}"
docker push "${SLUG}":"${VERSION}"
