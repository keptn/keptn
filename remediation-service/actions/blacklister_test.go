package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMixerTemplates(t *testing.T) {

	b := NewBlackLister()
	templates, err := b.getMixerTemplates("..")

	assert.Nil(t, err)
	assert.Equal(t, 3, len(templates))
	assert.NotEqual(t, "", templates[ipHandler])
	assert.NotEqual(t, "", templates[ipRule])
	assert.NotEqual(t, "", templates[ipInstance])
}

func TestAddIP(t *testing.T) {

	const emptyHandler = `apiVersion: config.istio.io/v1alpha2
kind: handler
metadata:
  name: blacklistip
spec:
  compiledAdapter: listchecker
  params:
    overrides: []
    blacklist: true
    entryType: IP_ADDRESSES
`
	const expectedHandler = `apiVersion: config.istio.io/v1alpha2
kind: handler
metadata:
  name: blacklistip
spec:
  compiledAdapter: listchecker
  params:
    blacklist: true
    entryType: IP_ADDRESSES
    overrides:
    - 127.0.0.1/32
`

	b := NewBlackLister()
	actualContent, err := b.addIP(emptyHandler, "127.0.0.1")

	assert.Nil(t, err)
	assert.Equal(t, expectedHandler, actualContent)
}
