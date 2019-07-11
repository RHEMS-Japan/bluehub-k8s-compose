#!/bin/bash

kubectl apply -f nginx/nginx-deployment.yaml
kubectl apply -f bluehub-gateway/bluehub-gateway-deployment.yaml
kubectl apply -f bluehub-router/bluehub-router-deployment.yaml
#kubectl apply -f private-registry/private-registry-deployment.yaml
kubectl apply -f nginx/nginx-service.yaml
kubectl apply -f bluehub-gateway/bluehub-gateway-service.yaml
kubectl apply -f bluehub-router/bluehub-router-service.yaml
#kubectl apply -f private-registry/private-registry-service.yaml
