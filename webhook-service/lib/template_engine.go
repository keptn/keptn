package lib

import (
	"bytes"
	"text/template"
)

//go:generate moq  -pkg fake -out ./fake/template_engine_mock.go . ITemplateEngine
type ITemplateEngine interface {
	ParseTemplate(data interface{}, templateStr string) (string, error)
}

type TemplateEngine struct{}

func (t *TemplateEngine) ParseTemplate(data interface{}, templateStr string) (string, error) {
	tmpl, err := template.New("").Parse(templateStr)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
