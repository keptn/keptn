#!/bin/bash

YLW='\033[1;33m'
NC='\033[0m'

CREDS=./creds.json
rm $CREDS 2> /dev/null

echo -e "${YLW}Please enter the credentials as requested below: ${NC}"
read -p "GitHub User Name: " GITU 
read -p "GitHub Personal Access Token: " GITAT
read -p "GitHub User Email: " GITE
read -p "GitHub Organization: " GITO
read -p "GKE Cluster Name: " CLN
read -p "GKE Cluster Zone: " CLZ
read -p "GKE Project: " PROJ
echo ""

echo ""
echo -e "${YLW}Please confirm all are correct: ${NC}"
echo "GitHub User Name: $GITU"
echo "GitHub Personal Access Token: $GITAT"
echo "GitHub User Email: $GITE"
echo "GitHub Organization: $GITO"
echo "GKE Cluster Name: $CLN"
echo "GKE Cluster Zone: $CLZ"
echo "GKE Project: $PROJ"
read -p "Is this all correct? (y/n) : " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]
then
    source ./defineCredentialsUtils.sh
    replaceCreds
fi

cat $CREDS
echo ""
echo "The credentials file can be found here:" $CREDS
echo ""

