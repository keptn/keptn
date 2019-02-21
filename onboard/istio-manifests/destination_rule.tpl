apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{ chartName }}-{{ serviceName }}-{{ environment }}-destination
spec:
  host: {{ chartName }}-{{ serviceName }}.{{ environment }}.svc.cluster.local
  subsets:
  - name: blue
    labels:
      deployment: {{ serviceName }}-blue
  - name: green
    labels:
      deployment: {{ serviceName }}-green