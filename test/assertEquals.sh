#!/bin/bash

# usage: ./assertEquals $origin_value $comparison_value

if [[ "$1" != "$2" ]];
then
    echo "not equal"
    exit 1
else
    echo "equal"
    exit 0
fi
