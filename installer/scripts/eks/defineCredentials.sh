#!/bin/bash

YLW='\033[1;33m'
NC='\033[0m'

echo -e "${YLW}Please enter the credentials as requested below: ${NC}"
read -p "Cluster Name: " CLN
read -p "AWS Region: " RG
echo ""

echo ""
echo -e "${YLW}Please confirm all are correct: ${NC}"
echo "Cluster Name: $CLN"
echo "AWS Region: $RG"
read -p "Is this all correct? (y/n) : " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]
then
    CREDS=./creds.json
    rm $CREDS 2> /dev/null
    cat ./aks/creds.sav | sed 's~CLUSTER_NAME_PLACEHOLDER~'"$CLN"'~' | \
      sed 's~AWS_REGION~'"$RG"'~' >> $CREDS

fi

cat $CREDS
echo ""
echo "The credentials file can be found here:" $CREDS
echo ""

