#!/bin/sh

REGISTRY_URI=$(kubectl describe svc docker-registry -n cicd | grep IP: | sed 's~IP:[ \t]*~~')
CHANNEL_URI=$(kubectl describe channel keptn-channel -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
NEW_ARTEFACT_CHANNEL=$(kubectl describe channel new-artefact -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
START_DEPLOYMENT_CHANNEL=$(kubectl describe channel start-deployment -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
DEPLOYMENT_FINISHED_CHANNEL=$(kubectl describe channel deployment-finished -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
START_TESTS_CHANNEL=$(kubectl describe channel start-tests -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
TESTS_FINISHED_CHANNEL=$(kubectl describe channel tests-finished -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
START_EVALUATION_CHANNEL=$(kubectl describe channel start-evaluation -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
EVALUATION_DONE_CHANNEL=$(kubectl describe channel evaluation-done -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')

rm -f config/gen/event-broker.yaml

cat config/event-broker.yaml | \
  sed 's~CHANNEL_URI_PLACEHOLDER~'"$CHANNEL_URI"'~' | \
  sed 's~NEW_ARTEFACT_CHANNEL_PLACEHOLDER~'"$NEW_ARTEFACT_CHANNEL"'~' | \
  sed 's~START_DEPLOYMENT_CHANNEL_PLACEHOLDER~'"$START_DEPLOYMENT_CHANNEL"'~' | \
  sed 's~DEPLOYMENT_FINISHED_CHANNEL_PLACEHOLDER~'"$DEPLOYMENT_FINISHED_CHANNEL"'~' | \
  sed 's~START_TESTS_CHANNEL_PLACEHOLDER~'"$START_TESTS_CHANNEL"'~' | \
  sed 's~TESTS_FINISHED_CHANNEL_PLACEHOLDER~'"$TESTS_FINISHED_CHANNEL"'~' | \
  sed 's~START_EVALUATION_CHANNEL_PLACEHOLDER~'"$START_EVALUATION_CHANNEL"'~' | \
  sed 's~EVALUATION_DONE_CHANNEL_PLACEHOLDER~'"$EVALUATION_DONE_CHANNEL"'~' | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/event-broker.yaml 
  
kubectl apply -f config/gen/event-broker.yaml