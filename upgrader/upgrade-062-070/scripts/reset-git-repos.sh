#!/bin/sh

cd /data/config/ || exit

for FILE in *; do
    if [ -d "$FILE" ]; then
        # shellcheck disable=SC2164
        cd "$FILE"
        git reset --hard
        # shellcheck disable=SC2164
        cd "$OLDPWD"
    fi
done
