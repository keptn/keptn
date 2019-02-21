apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ gitHubOrg }}-{{ serviceName }}-{{ environment }}
spec:
  hosts:
  - "{{ serviceName }}.{{ environment }}.{{ gitHubOrg }}.{{ ingressGatewayIP }}.xip.io"
  gateways:
  - {{ gitHubOrg }}-gateway
  http:
    - route:
      - destination:
          host: {{chartName}}-{{ serviceName }}.{{ environment }}.svc.cluster.local
          subset: blue
        weight: 0
      - destination:
          host: {{chartName}}-{{ serviceName }}.{{ environment }}.svc.cluster.local
          subset: green
        weight: 100