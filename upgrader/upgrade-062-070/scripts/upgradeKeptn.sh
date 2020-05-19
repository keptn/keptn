#!/bin/bash
source ./utils.sh

# Upgrade from Helm v2 to Helm v3
helm init --client-only
verify_install_step $? "Helm init failed."
RELEASES=$(helm list -aq)
verify_install_step $? "Helm list failed."
echo $RELEASES

helm3 plugin install https://github.com/helm/helm-2to3
verify_install_step $? "Helm-2to3 plugin installation failed."
yes y | helm3 2to3 move config
verify_install_step $? "Helm-2to3 move of config failed."

for release in $RELEASES; do
  helm3 2to3 convert $release --dry-run
  verify_install_step $? "Helm2-to3 release convertion dry-run failed"
  helm3 2to3 convert $release
  verify_install_step $? "Helm2-to3 release convertion failed"
done

yes y | helm3 2to3 cleanup --tiller-cleanup
verify_install_step $? "Helm2-to3 cleanup failed"

./upgrade-mongodb $MONGODB_URL $CONFIGURATION_SERVICE_URL

