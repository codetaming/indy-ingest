#!/usr/bin/env bash
docker run --env-file ./env -p 9000:9000 quay.io/codetaming/indy-ingest:latest