#!/bin/bash

YLW='\033[1;33m'
NC='\033[0m'

CREDS=./creds.json
rm $CREDS 2> /dev/null

echo -e "${YLW}Please enter the credentials as requested below: ${NC}"
read -p "Dynatrace Tenant (default=$DTENV): " DTENVC
read -p "Dynatrace API Token (default=$DTAPI): " DTAPIC
read -p "Dynatrace PaaS Token (default=$DTAPI): " DTPAAST
read -p "GitHub User Name: " GITU 
read -p "GitHub Personal Access Token: " GITAT
read -p "GitHub User Email: " GITE
read -p "GitHub Organization: " GITO
echo ""

if [[ $DTENV = '' ]]
then 
    DTENV=$DTENVC
fi

if [[ $DTAPI = '' ]]
then 
    DTAPI=$DTAPIC
fi

echo ""
echo -e "${YLW}Please confirm all are correct: ${NC}"
echo "Dynatrace Tenant: $DTENV"
echo "Dynatrace API Token: $DTAPI"
echo "Dynatrace PaaS Token: $DTPAAST"
echo "GitHub User Name: $GITU"
echo "GitHub Personal Access Token: $GITAT"
echo "GitHub User Email: $GITE"
echo "GitHub Organization: $GITO" 
read -p "Is this all correct? (y/n) : " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]
then
    rm $CREDS 2> /dev/null
    cat ./creds.sav | sed 's~DYNATRACE_TENANT_PLACEHOLDER~'"$DTENV"'~' | \
      sed 's~DYNATRACE_API_TOKEN~'"$DTAPI"'~' | \
      sed 's~DYNATRACE_PAAS_TOKEN~'"$DTPAAST"'~' | \
      sed 's~GITHUB_USER_NAME_PLACEHOLDER~'"$GITU"'~' | \
      sed 's~PERSONAL_ACCESS_TOKEN_PLACEHOLDER~'"$GITAT"'~' | \
      sed 's~GITHUB_USER_EMAIL_PLACEHOLDER~'"$GITE"'~' | \
      sed 's~GITHUB_ORG_PLACEHOLDER~'"$GITO"'~' >> $CREDS
fi

cat $CREDS
echo ""
echo "The credentials file can be found here:" $CREDS
echo ""