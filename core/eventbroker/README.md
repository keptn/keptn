# Event Broker

This component accepts incoming events from sources such as GitHub Webhooks, and pushes them into an internal knative eveinting queue. This queue will fan out incoming messages to subscribers.

## Installation

1. Create the channel:
```
kubectl apply -f config/channel.yaml
```
2. Retrieve the channel's URI, using the following command:

```
$ kubectl describe channel keptn-channel

Name:         keptn-channel
Namespace:    default
Labels:       <none>
Annotations:  kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"eventing.knative.dev/v1alpha1","kind":"Channel","metadata":{"annotations":{},"name":"keptn-channel","namespace":"default"},"spec":{"prov...
API Version:  eventing.knative.dev/v1alpha1
Kind:         Channel
Metadata:
  Creation Timestamp:  2019-02-08T08:37:47Z
  Finalizers:
    in-memory-channel-controller
  Generation:        15
  Resource Version:  912470
  Self Link:         /apis/eventing.knative.dev/v1alpha1/namespaces/default/channels/keptn-channel
  UID:               cf9756ec-2b7c-11e9-8026-42010a8000ae
Spec:
  Generation:  15
  Provisioner:
    API Version:  eventing.knative.dev/v1alpha1
    Kind:         ClusterChannelProvisioner
    Name:         in-memory-channel
  Subscribable:
    Subscribers:
      Ref:
        Name:          keptn-operator-subscription
        Namespace:     default
        UID:           748613f7-2de8-11e9-9906-42010a800049
      Subscriber URI:  http://keptn-operator.default.svc.cluster.local/
Status:
  Address:
    Hostname:  keptn-channel-channel-pq8xw.default.svc.cluster.local
  Conditions:
    Last Transition Time:  2019-02-08T08:37:47Z
    Severity:              Error
    Status:                True
    Type:                  Addressable
    Last Transition Time:  2019-02-08T08:37:47Z
    Severity:              Error
    Status:                True
    Type:                  Provisioned
    Last Transition Time:  2019-02-08T08:37:47Z
    Severity:              Error
    Status:                True
    Type:                  Ready
Events:                    <none>
```

3. Replace *CHANNEL_URI_PLACEHOLDER* in config/event-broker.yaml with the Hostname of the channel (in this example: *keptn-channel-channel-pq8xw.default.svc.cluster.local*).

4. Deploy the event broker:
```
kubectl apply -f config/event-broker.yaml
```

5. Verify the deployment:

```
$ kubectl get pods

NAME                                                  READY     STATUS            RESTARTS   AGE
event-broker-00002-deployment-5977798687-8c8c7        3/3       Running           0          12s
```