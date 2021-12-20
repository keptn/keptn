#!/bin/bash

for i in {$1..$2}
do
   keptn create service service-$i --project=multiple-replicas
   sleep 2
done