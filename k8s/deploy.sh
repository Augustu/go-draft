#!/bin/bash

echo "cp mission binary"
cp /home/augustu/Work/code/mission/mission /home/augustu/Work/code/mission-release/docker/

echo "remove k8s deployment"
cd /home/augustu/Work/code/mission-release/k8s
kubectl -n default delete -f mission.yaml

echo "remove old image"
docker image rm mission:v1

echo "build docker image"
cd /home/augustu/Work/code/mission-release/docker
docker build -t mission:v1 .

echo "deploy k8s image"
cd /home/augustu/Work/code/mission-release/k8s
kubectl -n default apply -f mission.yaml
