#!/bin/bash

# setup ssh certificate
USSH=$HOME/.ssh
mkdir -p $USSH
ssh-keygen -t rsa -N '' -f $USSH/ci_id_rsa
cat >> $USSH/config <<EOF
Host localhost
  StrictHostKeyChecking no
EOF

# Allow User CI key to login as root
sudo bash <<EOF
mkdir -p /root/.ssh
cat $USSH/ci_id_rsa.pub >> /root/.ssh/authorized_keys
chmod g-rw,o-rw /root/.ssh /root/.ssh/* $USSH/* $USSH
EOF

# make sure you have a profile is set correctly
minishift profile set keptn-dev
# minimum memory required for the minishift VM
minishift config set memory 4GB
# the minimum required vCpus for the minishift vm
minishift config set cpus 2
# Add new user called admin with password admin having role cluster-admin
minishift addons enable admin-user
# Allow the containers to be run with uid 0
minishift addons enable anyuid
minishift start --vm-driver=generic --remote-ipaddress 127.0.0.1 --remote-ssh-user root --remote-ssh-key $HOME/.ssh/ci_id_rsa

# Enable admission controller webhooks
# The configuration stanzas below look weird and are just to workaround for:
# https://bugzilla.redhat.com/show_bug.cgi?id=1635918
minishift openshift config set --target=kube --patch '{
    "admissionConfig": {
        "pluginConfig": {
            "ValidatingAdmissionWebhook": {
                "configuration": {
                    "apiVersion": "apiserver.config.k8s.io/v1alpha1",
                    "kind": "WebhookAdmission",
                    "kubeConfigFile": "/dev/null"
                }
            },
            "MutatingAdmissionWebhook": {
                "configuration": {
                    "apiVersion": "apiserver.config.k8s.io/v1alpha1",
                    "kind": "WebhookAdmission",
                    "kubeConfigFile": "/dev/null"
                }
            }
        }
    }
}'

# wait until the kube-apiserver is restarted
echo "Waiting for login ..."
until oc login -u admin -p admin; do sleep 5; done;

echo "Setting policies"

oc adm policy --as system:admin add-cluster-role-to-user cluster-admin admin
oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:default:default
oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:kube-system:default

# wait a little bit to ensure the cluster is ready
sleep 30
