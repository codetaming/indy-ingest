#!/usr/bin/env bash
docker run --env-file ./env -p 9000:9000 docker.io/codetaming/indy-ingest:latest