#!/bin/bash

YLW='\033[1;33m'
NC='\033[0m'

type hub &> /dev/null

if [ $? -ne 0 ]
then
    echo "Please install the 'hub' command: https://hub.github.com/"
    exit 1
fi

if [ -z $1 ]
then
    echo "Please provide the target GitHub orgainzation as parameter:"
    echo ""
    echo "  e.g.: ./forkGitHubRepositories.sh myorganization"
    echo ""
    exit 1
else
    ORG=$1
fi

HTTP_RESPONSE=`curl -s -o /dev/null -I -w "%{http_code}" https://github.com/$ORG`

if [ $HTTP_RESPONSE -ne 200 ]
then
    echo "GitHub organization doesn't exist - https://github.com/$ORG - HTTP status code $HTTP_RESPONSE"
    exit 1
fi

declare -a repositories=("carts" "catalogue" "front-end" "jenkins-release-branch" "k8s-deploy-production" "k8s-deploy-staging" "orders" "payment" "queue-master" "shipping" "sockshop-infrastructure" "user")

mkdir ../repositories
cd ../repositories

for repo in "${repositories[@]}"
do
    #FOLDER=$(echo $repo | cut -d '/' -f 5)
    echo -e "${YLW}Cloning https://github.com/keptn-sockshop/$repo ${NC}"
    git clone -q "https://github.com/keptn-sockshop/$repo"
    cd $repo
    echo -e "${YLW}Forking $repo to $ORG ${NC}"
    hub fork --org=$ORG
    cd ..
    echo -e "${YLW}Done. ${NC}"
done

cd ..
rm -rf repositories
mkdir repositories
cd repositories

for repo in "${repositories[@]}"
do
    TARGET_REPO="http://github.com/$ORG/$repo"
    echo -e "${YLW}Cloning $TARGET_REPO ${NC}"
    git clone -q $TARGET_REPO
    echo -e "${YLW}Done. ${NC}"
done
