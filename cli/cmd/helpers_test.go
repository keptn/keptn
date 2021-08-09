package cmd

import (
	"bytes"
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeWatcher struct {
}

type HelpersTestStruct struct {
	Value string `json:"value"`
}

func (fw fakeWatcher) Watch(ctx context.Context) (<-chan []*models.KeptnContextExtendedCE, context.CancelFunc) {
	ch := make(chan []*models.KeptnContextExtendedCE, 2)
	defer close(ch)
	ch <- []*models.KeptnContextExtendedCE{&models.KeptnContextExtendedCE{ID: "ID1"}}
	return ch, func() {}
}

func TestPrintEvents(t *testing.T) {

	aStruct := HelpersTestStruct{Value: "value"}
	var tests = []struct {
		buff    bytes.Buffer
		format  string
		content interface{}
		exp     string
	}{
		{
			bytes.Buffer{},
			"json",
			aStruct,
			`{"value":"value"}`,
		},
		{
			bytes.Buffer{},
			"",
			aStruct,
			`{"value":"value"}`,
		},
		{
			bytes.Buffer{},
			"yaml",
			aStruct,
			`value:value`,
		}, {
			bytes.Buffer{},
			"",
			HelpersTestStruct{Value: "<=+75%"},
			`{"value":"\u003c=+75%"}`,
		},
	}

	for _, e := range tests {
		PrintEvents(&e.buff, e.format, e.content)
		act := fullTrim(e.buff.String())
		if act != e.exp {
			t.Errorf("Print Events output: %s, expected: %s", act, e.exp)
		}
	}
}

func TestPrintEventWatcher(t *testing.T) {

	exp := `{
    "data": null,
    "id": "ID1",
    "source": null,
    "time": "0001-01-01T00:00:00Z",
    "type": null
}`

	fakeWatcher := fakeWatcher{}
	var buff bytes.Buffer
	PrintEventWatcher(context.TODO(), fakeWatcher, "json", &buff)
	assert.Equal(t, fullTrim(exp), fullTrim(buff.String()))
}

func fullTrim(str string) string {
	s := strings.ReplaceAll(str, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

func stringp(s string) *string {
	return &s
}

func boolp(b bool) *bool {
	return &b
}

func getTestAPI() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			if r.Method == http.MethodPost {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(``))
				return
			} else if r.Method == http.MethodGet {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(`{"user":"user", "password":"password"}`))
				return
			}
		}),
	)
	return ts
}
