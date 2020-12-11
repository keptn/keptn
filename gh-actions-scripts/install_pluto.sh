#!/bin/bash

# download
wget https://github.com/FairwindsOps/pluto/releases/download/v3.5.1/pluto_3.5.1_linux_amd64.tar.gz -O pluto.tar.gz
tar -zxvf pluto.tar.gz
if [ ! -d "$GITHUB_WORKSPACE/bin" ]; then
  mkdir "$GITHUB_WORKSPACE/bin"
fi

