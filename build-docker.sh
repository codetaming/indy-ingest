#!/usr/bin/env bash
make
docker build -f Dockerfile-arm -t codetaming/ingest-arm .
docker build -f Dockerfile -t codetaming/ingest .
docker login
docker push codetaming/ingest-arm
docker push codetaming/ingest