#!/bin/bash

echo "Removing sign-off messages from changelog..."
for file in {CHANGELOG,RELEASE-BODY}.md; do
  if [ -f "$file" ]; then
    echo "Replacing content in $file"
    # Reference: https://stackoverflow.com/a/1252191
    sed -e ':a' -e 'N' -e '$!ba' -e 's/\nSigned-off-by: .* <.*@.*>\n/ /g' "$file" > tmp
    mv tmp "$file"
  else
    echo "Not replacing anything since $file does not exist."
  fi
done
