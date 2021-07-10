#!/bin/bash


echo "remove k8s deployment"
kubectl -n dev delete -f httptest.yaml

echo "build httptest binary"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
strip main

echo "remove old image"
docker image rm httptest:v1

echo "build docker image"
docker build -t httptest:v1 .

echo "deploy k8s image"
kubectl -n dev apply -f httptest.yaml
