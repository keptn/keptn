package importer

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keptn/keptn/api/test/utils"
)

func TestRenderContent(t *testing.T) {

	type input struct {
		template string
		context  any
	}
	type expectation struct {
		renderedText string
	}

	tests := []struct {
		name         string
		inputs       input
		expectations expectation
	}{
		{
			name:         "Empty template",
			inputs:       input{template: "", context: nil},
			expectations: expectation{renderedText: ""},
		},
		{
			name: "Simple template with [[ ... ]] delimiters",
			inputs: input{
				template: "This is [[ .rendered ]], while this is not {{ .rendered }}",
				context:  map[string]string{"rendered": "awesomely rendered"},
			},
			expectations: expectation{renderedText: "This is awesomely rendered, while this is not {{ .rendered }}"},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				readCloser := io.NopCloser(strings.NewReader(tt.inputs.template))
				r, err := RenderContent(readCloser, tt.inputs.context)
				// Assert that whatever RenderContent returns implements io.ReadCloser (
				// it has to be substituted to the raw content io.ReadCloser passed down to the executors)
				assert.Implements(t, (*io.ReadCloser)(nil), r)
				assert.NoError(t, err)
				bytes, err := io.ReadAll(r)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectations.renderedText, string(bytes))
			},
		)
	}
}

func TestErrorRenderContent(t *testing.T) {

	type input struct {
		template io.ReadCloser
		context  any
	}
	type expectation struct {
		err           error
		errorContains string
	}

	tests := []struct {
		name         string
		inputs       input
		expectations expectation
	}{
		{
			name: "Broken reader",
			inputs: input{
				template: io.NopCloser(utils.NewTestReader([]byte("some text"), 0, true)),
				context:  nil,
			},
			expectations: expectation{
				errorContains: "error reading template: ",
			},
		},
		{
			name: "Invalid template - unclosed action",
			inputs: input{
				template: io.NopCloser(strings.NewReader("invalid template syntax [[ .unclosed_thing")),
				context:  nil,
			},
			expectations: expectation{
				errorContains: "error parsing template: ",
			},
		},
		{
			name: "Invalid template - non-exiting function",
			inputs: input{
				template: io.NopCloser(strings.NewReader("invalid template syntax [[ non-existing-function ]]")),
				context:  nil,
			},
			expectations: expectation{
				errorContains: "error parsing template: ",
			},
		},
		{
			name: "Error rendering template",
			inputs: input{
				template: io.NopCloser(strings.NewReader("trying to index non-existing field [[ index .field1 0 ]]")),
				context:  nil,
			},
			expectations: expectation{
				errorContains: "error rendering template: ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				r, err := RenderContent(tt.inputs.template, tt.inputs.context)
				assert.Nil(t, r)
				assert.Error(t, err)
				if tt.expectations.err != nil {
					assert.ErrorIs(t, err, tt.expectations.err)
				}
				if tt.expectations.errorContains != "" {
					assert.ErrorContains(t, err, tt.expectations.errorContains)
				}
			},
		)
	}
}
