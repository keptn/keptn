#!/bin/sh

cd /data/config/
for FILE in *; do
    if [ -d "$FILE" ]; then
        cd "$FILE"
        git reset --hard
        cd ..
    fi
done
