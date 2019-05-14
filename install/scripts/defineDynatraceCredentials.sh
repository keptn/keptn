#!/bin/bash

YLW='\033[1;33m'
NC='\033[0m'

CREDS=./creds_dt.json
rm $CREDS 2> /dev/null

echo -e "${YLW}Please enter the credentials as requested below: ${NC}"
read -p "Dynatrace Tenant {your-domain}/e/{your-environment-id} for managed or {your-environment-id}.live.dynatrace.com for SaaS (default=$DTENV): " DTENVC
read -p "Dynatrace API Token (default=$DTAPI): " DTAPIC
read -p "Dynatrace PaaS Token (default=$DTAPI): " DTPAAST
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
read -p "Is this all correct? (y/n) : " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]
then
    rm $CREDS 2> /dev/null
    cat ./creds_dt.sav | sed 's~DYNATRACE_TENANT_PLACEHOLDER~'"$DTENV"'~' | \
      sed 's~DYNATRACE_API_TOKEN~'"$DTAPI"'~' | \
      sed 's~DYNATRACE_PAAS_TOKEN~'"$DTPAAST"'~' >> $CREDS
fi

cat $CREDS
echo ""
echo "The credentials file can be found here:" $CREDS
echo ""
