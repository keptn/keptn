#!/bin/bash

if ! command -v npm &> /dev/null
then
    echo "npm could not be found, install npm first"
    exit
fi

cd ./bridge/

# install the dependencies
npm install

# accept the API URL and API TOKEN as the user inputs
read -p 'Enter API URL: ' API_URL
read -sp 'Enter API Token: ' API_TOKEN

# start the server
npm run start:dev
