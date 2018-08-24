#!/usr/bin/env bash
kubectx pi-k8s
kubens go
kubectl apply -f k8s/.
kubectl get all