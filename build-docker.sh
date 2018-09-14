#!/usr/bin/env bash
docker build -f Dockerfile-arm -t codetaming/indy-ingest-arm .
docker build -f Dockerfile -t codetaming/indy-ingest .
docker login
docker push codetaming/indy-ingest-arm
docker push codetaming/indy-ingest