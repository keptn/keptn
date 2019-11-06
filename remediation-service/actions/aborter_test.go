package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const vsWithAbortFault = `apiVersion: networking.istio.io/v1alpha3
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
      abort:
        httpStatus: 403
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

func TestAddAbortDelay(t *testing.T) {

	a := NewAborter()
	newVs, err := a.addAbort(virtualService, "2.2.2.2")

	assert.Nil(t, err)
	assert.Equal(t, vsWithAbortFault, newVs)
}

func TestValidRemoveAbort(t *testing.T) {

	a := NewAborter()
	newVs, err := a.removeAbort(vsWithAbortFault, "2.2.2.2")

	assert.Nil(t, err)
	assert.Equal(t, virtualService, newVs)
}

func TestInvalidRemoveAbort1(t *testing.T) {

	a := NewAborter()
	newVs, err := a.removeAbort(vsWithAbortFault, "1.1.1.1")

	assert.Nil(t, err)
	assert.Equal(t, vsWithAbortFault, newVs)
}

func TestInvalidRemoveAbort2(t *testing.T) {

	a := NewAborter()
	newVs, err := a.removeAbort(virtualService, "1.1.1.1")

	assert.Nil(t, err)
	assert.Equal(t, virtualService, newVs)
}

func TestValidDoubleRemoveAbort(t *testing.T) {

	const vsWithDoubleAbortFault = `apiVersion: networking.istio.io/v1alpha3
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
      abort:
        httpStatus: 403
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
  - fault:
      abort:
        httpStatus: 403
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
	a := NewAborter()
	newVs, err := a.removeAbort(vsWithDoubleAbortFault, "2.2.2.2")

	assert.Nil(t, err)
	assert.Equal(t, virtualService, newVs)
}
