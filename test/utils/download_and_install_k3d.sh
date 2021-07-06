#!/bin/bash

# install k3d version from github (see https://github.com/rancher/k3s/releases)
K3D_VERSION=${K3D_VERSION:-"v4.4.4"}
# install K3d
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh --max-time 300 | TAG=${K3D_VERSION} bash
