#!/bin/bash

# download
wget https://github.com/FairwindsOps/pluto/releases/download/v3.5.1/pluto_3.5.1_linux_amd64.tar.gz -O pluto.tar.gz
tar -zxvf pluto.tar.gz
if [ ! -d "$HOME/bin" ]; then
  mkdir "$HOME/bin"
fi

mv pluto "$HOME/bin"
export PATH=$PATH:"$HOME/bin"
