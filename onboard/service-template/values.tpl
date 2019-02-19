replicaCount: 1
image:
    repository: null
    tag: null
    pullPolicy: IfNotPresent
service:
    name: {{ microServiceName }}
    type: LoadBalancer
    externalPort: 8080
    internalPort: 8080
resources:
    limits:
        cpu: 100m
        memory: 128Mi
    requests:
        cpu: 100m
        memory: 128Mi