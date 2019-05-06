#!/bin/bash
kubectl delete -f config/event-broker-ext.yaml --ignore-not-found
kubectl apply -f config/event-broker-ext.yaml