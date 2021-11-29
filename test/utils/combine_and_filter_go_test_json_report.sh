#!/bin/bash

SOURCE_TEST_REPORT_FOLDER=$1
DESTINATION_TEST_REPORT_FILE=$2
REWRITE_TEST_TARGET=$3

if [ $# -ne 3 ]; then
  echo "Usage: $0 SOURCE_TEST_REPORT_FOLDER DESTINATION_TEST_REPORT_FILE REWRITE_TEST_TARGET"
  exit
fi

if [ -z "$SOURCE_TEST_REPORT_FOLDER" ]; then
  echo "No source report folder set, exiting..."
  exit 1
fi

if [ -z "$DESTINATION_TEST_REPORT_FILE" ]; then
  echo "No destination report filename, exiting..."
  exit 1
fi

if [ -z "$REWRITE_TEST_TARGET" ]; then
  echo "No test target for rewriting set, exiting..."
  exit 1
fi

echo "Combining test report files from source folder..."
# shellcheck disable=SC2066
for file in "$SOURCE_TEST_REPORT_FOLDER"/*; do
  echo "Appending $file ..."
  cat "$file" >> "$DESTINATION_TEST_REPORT_FILE"
done

echo "Filtering combined test report file..."
# shellcheck disable=SC2002
filtered_report=$(cat "$DESTINATION_TEST_REPORT_FILE" | jq -c "select( (.Action == \"pass\" or .Action == \"fail\") and .Test != null ) | {\"test\": .Test, \"$REWRITE_TEST_TARGET\": .Action}")
echo "$filtered_report" > "$DESTINATION_TEST_REPORT_FILE"
