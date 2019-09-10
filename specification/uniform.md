# Uniform (DRAFT)

A `uniform` describes the services that listen to Keptn channels for events and defines to which service subscribes to which channel. Configuration parameters can be passed either as environment variable or references to config maps.

Ultimately, the `uniform` will be a custom Kubernetes resource, as suggested below. The `uniform` controller then will take care of managing subscribed services and subscriptions. This enables Keptn users to change their continuous delivery services and version with ease, because the `uniform` is the one central location where this information is stored.

```yaml
---
apiVersion: sh.keptn/v1alpha1
kind: Uniform
metadata:
  name: bookinfo-uniform
  namespace: keptn
spec:
  services:
  - name: slack-trail
    image: keptn/slack-service:0.1.2
    env:
    - name: SLACK_WEBHOOK
      value: "https://hooks.slack.com/services/TXXXXXXXX/BXXXXXXX/WXXXXXXXXXXXXXXXXX"
    subscribedchannels:
    - new-artifact
    - configuration-changed
    - deployment-done
    - tests-done
    - evaluation-done
    - problem
  - name: gitops-operator
    image: keptn/github-service:0.3.4
    env:
    - name: GITHUB_USER
      value: "johndoe"
    - name: GITHUB_TOKEN
      valueFrom:
        secretKeyRef:
            name: github-token
            key: GITHUB_TOKEN
    subscribedchannels:
    - new-artifact
  - name: cd-operator
    image: keptn/jenkins-service:0.1.1
    subscribedchannels:
    - configuration-changed
    - deployment-done
    - tests-done
    - evaluation-done
  - name: automated-operations
    image: keptn/servicenow-service:0.2.1
    subscribedchannels:
    - problem
```