package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDelay(t *testing.T) {

	const expectedVirtualService = `apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  creationTimestamp: null
  name: carts
spec:
  gateways:
  - sockshop-production-gateway.sockshop-production
  - mesh
  hosts:
  - carts.sockshop-production.1.1.1.1.xip.io
  - carts
  http:
  - fault:
      delay:
        fixedDelay: 5s
        percent: 100
    match:
    - headers:
        X-Forwarded-For:
          exact: 2.2.2.2
    route:
    - destination:
        host: carts-canary.sockshop-production.svc.cluster.local
    - destination:
        host: carts-primary.sockshop-production.svc.cluster.local
      weight: 100
  - route:
    - destination:
        host: carts-canary.sockshop-production.svc.cluster.local
    - destination:
        host: carts-primary.sockshop-production.svc.cluster.local
      weight: 100
`

	s := NewSlower()
	newVs, err := s.addDelay(virtualService, "2.2.2.2", "5s")

	assert.Nil(t, err)
	assert.Equal(t, expectedVirtualService, newVs)
}
