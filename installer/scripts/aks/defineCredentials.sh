#!/bin/bash

YLW='\033[1;33m'
NC='\033[0m'

echo -e "${YLW}Please enter the credentials as requested below: ${NC}"
read -p "GitHub User Name: " GITU 
read -p "GitHub Personal Access Token: " GITAT
read -p "GitHub Organization: " GITO
read -p "Cluster Name: " CLN
read -p "Azure Resource Group: " RG
read -p "Azure Subscription: " AS
echo ""

echo ""
echo -e "${YLW}Please confirm all are correct: ${NC}"
echo "GitHub User Name: $GITU"
echo "GitHub Personal Access Token: $GITAT"
echo "GitHub Organization: $GITO"
echo "Cluster Name: $CLN"
echo "Azure Resource Group: $RG"
echo "Azure Subscription: $AS"
read -p "Is this all correct? (y/n) : " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]
then
    CREDS=./creds.json
    rm $CREDS 2> /dev/null
    cat ./aks/creds.sav | sed 's~GITHUB_USER_NAME_PLACEHOLDER~'"$GITU"'~' | \
      sed 's~PERSONAL_ACCESS_TOKEN_PLACEHOLDER~'"$GITAT"'~' | \
      sed 's~CLUSTER_NAME_PLACEHOLDER~'"$CLN"'~' | \
      sed 's~AZURE_RESOURCE_GROUP~'"$RG"'~' | \
      sed 's~AZURE_SUBSCRIPTION~'"$AS"'~' | \
      sed 's~GITHUB_ORG_PLACEHOLDER~'"$GITO"'~' >> $CREDS

fi

cat $CREDS
echo ""
echo "The credentials file can be found here:" $CREDS
echo ""

