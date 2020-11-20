#!/bin/bash

source test/utils.sh

function cleanup() {
  sleep 2
  echo "Deleting Project ${PROJECT}"
  keptn delete project ${PROJECT}

  echo "Deleting echo-service deployment"
  kubectl delete deployments/echo-service -n keptn

  echo "Deleting echo-service service2"
  kubectl delete services/echo-service -n keptn

  echo "<END>"
}

function verify_event_not_null() {
  if [[ $1 == "null" ]]; then
    return -1
  fi
}

trap cleanup EXIT

# get keptn API details
KEPTN_ENDPOINT=http://$(kubectl -n keptn get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/api
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)

#test configuration
PROJECT="linking-stages-project"
SERVICE="linking-stages-service"

ECHO_SVC_IMG="docker.io/warber/keptnsandbox_echo-service"

echo "Installing echo service"
kubectl apply -f - <<EOF
---
# Deployment of our echo-service
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-service
  namespace: keptn
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: echo-service
      app.kubernetes.io/instance: keptn
  replicas: 1
  template:
    metadata:
      labels:
        app.kubernetes.io/name: echo-service
        app.kubernetes.io/instance: keptn
        app.kubernetes.io/part-of: keptn-keptn
        app.kubernetes.io/component: control-plane
        app.kubernetes.io/version: develop
    spec:
      containers:
        - name: echo-service
          image: warber/keptnsandbox_echo-service:72478fe-dirty #keptnsandbox/echo-service:latest # Todo: Replace this with your image name
          ports:
            - containerPort: 8080
          env:
            - name: EVENTBROKER
              value: 'http://localhost:8081/event'
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
        - name: distributor
          image: keptn/distributor:latest
          livenessProbe:
            httpGet:
              path: /health
              port: 10999
            initialDelaySeconds: 5
            periodSeconds: 5
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
          resources:
            requests:
              memory: "32Mi"
              cpu: "50m"
            limits:
              memory: "128Mi"
              cpu: "500m"
          env:
            - name: PUBSUB_URL
              value: 'nats://keptn-nats-cluster'
            - name: PUBSUB_TOPIC
              value: 'sh.keptn.>'
            - name: PUBSUB_RECIPIENT
              value: '127.0.0.1'
            - name: PUBSUB_RECIPIENT_PATH
              value: '/v1/event'
---
# Expose echo-service via Port 8080 within the cluster
apiVersion: v1
kind: Service
metadata:
  name: echo-service
  namespace: keptn
  labels:
    app.kubernetes.io/name: shipyard-controller
    app.kubernetes.io/instance: keptn
    app.kubernetes.io/part-of: keptn-keptn
    app.kubernetes.io/component: control-plane
spec:
  ports:
    - port: 8080
      targetPort: 8080
      protocol: TCP
  selector:
    app.kubernetes.io/name: echo-service
    app.kubernetes.io/instance: keptn
EOF

echo "Testing link staging..."

echo "Creating a new project without Git upstream"
keptn create project ${PROJECT} --shipyard=./test/assets/linking_stages_shipyard.yaml
sleep 1

echo "Creating a new service"
keptn create service ${SERVICE} --project ${PROJECT}
sleep 1

echo "Sending trigger echosequence event event"
keptn_context_id=$(send_event_json ./test/assets/trigger_echosequence_event.json)
sleep 20


verify_event_not_null $(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.firststage.echosequence.triggered $KEPTN_ENDPOINT $KEPTN_API_TOKEN)
if [ "$?" -eq "-1" ];then
echo "Event for triggering first stage could not be fetched"
exit 2
fi 

verify_event_not_null $(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.firststage.echosequence.finished $KEPTN_ENDPOINT $KEPTN_API_TOKEN)
if [ "$?" -eq "-1" ];then
echo "Event for finishing first stage could not be fetched"
exit 2
fi 

verify_event_not_null $(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.secondstage.echosequence.triggered $KEPTN_ENDPOINT $KEPTN_API_TOKEN)
if [ "$?" -eq "-1" ];then
echo "Event for triggering second stage could not be fetched"
exit 2
fi 

verify_event_not_null $(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.secondstage.echosequence.finished $KEPTN_ENDPOINT $KEPTN_API_TOKEN)
if [ "$?" -eq "-1" ];then
echo "Event for finishing second stage could not be fetched"
exit 2
fi 