#!/usr/bin/env bash
docker build -f Dockerfile-arm -t codetaming/indy-ingest-arm .
docker build -f Dockerfile -t codetaming/indy-ingest .
docker login
docker push codetaming/indy-ingest-arm
docker push codetaming/indy-ingest
docker build -f Dockerfile-arm -t docker-registry.k8s.codetaming.org/indy-ingest-arm/indy-ingest-arm .
docker push docker-registry.k8s.codetaming.org/indy-ingest-arm