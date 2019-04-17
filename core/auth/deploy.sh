#!/bin/sh
kubectl delete -f config/authenticator.yaml --ignore-not-found
kubectl apply -f config/authenticator.yaml