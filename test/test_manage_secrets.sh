#!/bin/bash

# shellcheck disable=SC1091

source test/utils.sh

function cleanup() {
  echo "<END>"
  return 0
}

trap cleanup EXIT

KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}

echo "########################################"
echo "TEST-1: Creating and deleting secrets"
echo "########################################"

SECRET_1="my-new-secret"
SECRET_2="my-new-secret-2"

keptn create secret $SECRET_1 --from-literal="mykey1=myvalue1"
verify_test_step $? "Failed to create secret $SECRET_1"

keptn create secret $SECRET_2 --from-literal="mykey2=myvalue2"
verify_test_step $? "Failed to create secret $SECRET_2"


kubectl get secrets $SECRET_1 -n "$KEPTN_NAMESPACE"
verify_test_step $? "Secret $SECRET_1 was not created"

kubectl get secrets $SECRET_2 -n "$KEPTN_NAMESPACE"
verify_test_step $? "Secret $SECRET_2 was not created"

roles_response=$(kubectl get roles keptn-secrets-default-read -n "$KEPTN_NAMESPACE" -ojson)
verify_using_jq "$roles_response" ".rules[0].resourceNames | length" 2
verify_using_jq "$roles_response" ".rules[0].resourceNames[0]" $SECRET_1
verify_using_jq "$roles_response" ".rules[0].resourceNames[1]" $SECRET_2

echo "Deleting a secret"
keptn delete secret $SECRET_1
verify_test_step $? "Failed to delete secret $SECRET_1"

keptn delete secret $SECRET_2
verify_test_step $? "Failed to delete secret $SECRET_2"

roles_response=$(kubectl get roles keptn-secrets-default-read -n "$KEPTN_NAMESPACE" -ojson)
verify_using_jq "$roles_response" ".rules[0].resourceNames | length" 0

echo "########################################"
echo "TEST-2: Creating and updating a secret"
echo "########################################"

SECRET_1="my-new-secret"
keptn create secret $SECRET_1 --from-literal="mykey1=myvalue1"
verify_test_step $? "Failed to create secret $SECRET_1"

get_secret_response=$(kubectl get secret $SECRET_1 -n "$KEPTN_NAMESPACE" -ojson)
old_secret_val=$(echo "$get_secret_response" | jq '.data.mykey1')

keptn update secret $SECRET_1 --from-literal="mykey1=changed-value"
verify_test_step $? "Failed to update secret $SECRET_1"

get_secret_response=$(kubectl get secret $SECRET_1 -n "$KEPTN_NAMESPACE" -ojson)
updated_secret_val=$(echo "$get_secret_response" | jq '.data.mykey1')
verify_not_equal "$old_secret_val" "$updated_secret_val"

keptn delete secret $SECRET_1
verify_test_step $? "Failed to delete secret $SECRET_1"

