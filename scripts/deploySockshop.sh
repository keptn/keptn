#!/bin/bash

cd ../repositories/sockshop-infrastructure
kubectl apply -f carts-db.yml
kubectl apply -f catalogue-db.yml
kubectl apply -f orders-db.yml
kubectl apply -f rabbitmq.yaml
kubectl apply -f user-db.yaml

# Apply services
declare -a repositories=("carts" "catalogue" "front-end" "orders" "payment" "queue-master" "shipping" "user")

for repo in "${repositories[@]}"
do
    cd ../repositories/$repo/manifest
    # Deploy service to dev
    kubectl apply -f ./$repo.yml

    # Deploy service to staging 
    cat $repo.yml | sed 's#namespace: .*#namespace: staging#' >> staging_tmp.yml
    kubectl apply -f ./staging_tmp.yml
    rm staging_tmp.yml

    # Deploy service to production
    cat $repo.yml | sed 's#namespace: .*#namespace: production#' >> production_tmp.yml
    kubectl apply -f ./production_tmp.yml
    rm production_tmp.yml
done