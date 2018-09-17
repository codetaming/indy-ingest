#!/usr/bin/env bash
kubectx pi-k8s
kubens indy-ingest
kubectl apply -f k8s/.
kubectl get all