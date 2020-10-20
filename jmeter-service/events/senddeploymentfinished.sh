#!/bin/bash

PROJECT=$1
STAGE=$2
SERVICE=$3
TESTSTRATEGY=$4
URL=${5//\//\\/}
VERSION=$6
USER=$7

if [[ -z "$VERSION" ]]; then
  VERSION="1.0"
fi

if [[ -z "$USER" ]]; then
  USER="noname"
fi

if [[ -z "$PROJECT" || -z "$STAGE" || -z "$SERVICE" || -z "$URL" || -z "$TESTSTRATEGY" ]]; then
  echo "Please specify project, stage, service, test strategy and URL. Optionally specify your version and user:"
  echo "Usage: ./senddeployfinished.sh PROJECT STAGE SERVICE TESTSTRATEGY URL [VERSION] [USER]"
  echo "Example: ./senddeployfinished.sh perfservice performance simplenodeservice performance 1.1 Andi"
  exit 1
fi

# Generate a temp file with replaced placeholders
inputfile="deployment.finished.event.placeholders.json"
tmpfile="deployment.finished.event.tmp.json"

if [ -f $tmpfile ] ; then
    rm -f $tmpfile
fi

sed -e "s/\$PROJECT/$PROJECT/" -e "s/\$STAGE/$STAGE/" -e "s/\$SERVICE/$SERVICE/" -e "s/\$TESTSTRATEGY/$TESTSTRATEGY/" -e "s/\$USER/$USER/" -e "s/\$VERSION/$VERSION/" -e "s/\$URL/$URL/" $inputfile >> $tmpfile

# now lets execute the keptn command
keptn send event --file=$tmpfile
# keptn send event --file $tmpfile