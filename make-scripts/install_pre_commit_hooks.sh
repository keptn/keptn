#!/bin/bash

curl https://pre-commit.com/install-local.py | python -

# pre-commit binary will be installed in $HOME/bin
export PATH=$PATH:"$HOME/bin"
