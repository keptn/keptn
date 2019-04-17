#!/bin/sh
kubectl delete -f config/event-broker.yaml --ignore-not-found
kubectl apply -f config/event-broker.yaml
