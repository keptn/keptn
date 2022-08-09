package importer

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

type templateRenderer struct{}

func (tr *templateRenderer) RenderContent(raw io.ReadCloser, context any) (io.ReadCloser, error) {
	defer raw.Close()
	bytes, err := io.ReadAll(raw)

	if err != nil {
		return nil, fmt.Errorf("error reading template: %w", err)
	}

	rendered, err := tr.RenderString(string(bytes), context)

	if err != nil {
		return nil, err
	}
	return io.NopCloser(strings.NewReader(rendered)), nil
}
func (tr *templateRenderer) RenderString(raw string, context any) (string, error) {

	t, err := template.New("template").
		Delims("[[", "]]").
		Option("missingkey=error").
		Parse(raw)

	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	buf := new(strings.Builder)
	err = t.Execute(buf, context)
	if err != nil {
		return "", fmt.Errorf("error rendering template: %w", err)
	}
	return buf.String(), nil
}
