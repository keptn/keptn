package importer

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

func RenderContent(raw io.ReadCloser, context any) (io.ReadCloser, error) {
	defer raw.Close()
	bytes, err := io.ReadAll(raw)

	if err != nil {
		return nil, fmt.Errorf("error reading template: %w", err)
	}

	t, err := template.New("template").
		Delims("[[", "]]").
		Option("missingkey=error").
		Parse(string(bytes))

	if err != nil {
		return nil, fmt.Errorf("error parsing template: %w", err)
	}

	buf := new(strings.Builder)
	err = t.Execute(buf, context)
	if err != nil {
		return nil, fmt.Errorf("error rendering template: %w", err)
	}
	return io.NopCloser(strings.NewReader(buf.String())), nil
}
