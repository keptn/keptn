#!/bin/bash

# check if npm command is present
if ! command -v npm &> /dev/null
then
    echo "npm could not be found, install npm first"
    exit
fi

cd ./bridge/

# Check if node_modules is absent
if [ -d ./bridge/node_modules ]; then
    # install the dependencies
    npm install
fi

if [[ -z $API_URL && -z $API_TOKEN ]]; then
    # accept the API URL and API TOKEN as the user inputs
    read -p 'Enter API URL: ' API_URL
    read -sp 'Enter API Token: ' API_TOKEN

    if [ -z $API_URL ]; then
        echo "Enter valid API URL"
        exit 0
    fi

    if [ -z $API_TOKEN ]; then
        echo "API Token is left blank, it will be automatically pulled via kubectl."
        exit 0
    fi
elif [[ -z $API_URL && ! -z $API_TOKEN ]]; then
    read -p 'Enter API URL: ' API_URL
else
    TOKENLEN = ${#API_TOKEN}
    s = $(printf "%-${TOKENLEN}s" "*")
    echo "API URL and API Token already set."
    echo "API URL:" $API_URL
    echo "API Token:" "${s// /*}"
fi

# start the server
npm run start:dev
