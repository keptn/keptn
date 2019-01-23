#!/bin/bash

cd ../repositories/sockshop-infrastructure
kubectl apply -f manifests/carts-db.yaml
kubectl apply -f manifests/catalogue-db.yaml
kubectl apply -f manifests/orders-db.yaml
kubectl apply -f manifests/rabbitmq.yaml
kubectl apply -f manifests/user-db.yaml

cd ..

# Apply services
declare -a repositories=("carts" "catalogue" "front-end" "orders" "payment" "queue-master" "shipping" "user")

for repo in "${repositories[@]}"
do
    cd $repo/manifest
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

    cd ../..
done
cd ../scripts