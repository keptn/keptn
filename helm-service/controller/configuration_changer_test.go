package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
)

func TestApplyFileChanges(t *testing.T) {

	meta := &chart.Metadata{
		Name: "test.chart",
	}
	const templateFileName = "templates/test.yaml"
	originalContent := []byte("originalContent")
	newContent := []byte("newContent")

	inputChart := chart.Chart{Metadata: meta}
	template := &chart.File{Name: templateFileName, Data: originalContent}
	inputChart.Templates = append(inputChart.Templates, template)

	fileChanges := make(map[string]string)
	fileChanges[templateFileName] = string(newContent)
	applyFileChanges(fileChanges, &inputChart)
	assert.Equal(t, templateFileName, inputChart.Templates[0].Name)
	assert.Equal(t, newContent, inputChart.Templates[0].Data)
}

func TestAddFile(t *testing.T) {

	meta := &chart.Metadata{
		Name: "test.chart",
	}
	const templateFileName = "templates/test.yaml"
	originalContent := []byte("originalContent")
	newContent := []byte("newContent")

	inputChart := chart.Chart{Metadata: meta}
	template := &chart.File{Name: templateFileName, Data: originalContent}
	inputChart.Templates = append(inputChart.Templates, template)

	fileChanges := make(map[string]string)
	const newTemplateFileName = "templates/newFile.yaml"
	fileChanges["templates/newFile.yaml"] = string(newContent)
	applyFileChanges(fileChanges, &inputChart)
	assert.Equal(t, 2, len(inputChart.Templates))
	assert.Equal(t, templateFileName, inputChart.Templates[0].Name)
	assert.Equal(t, originalContent, inputChart.Templates[0].Data)
	assert.Equal(t, newTemplateFileName, inputChart.Templates[1].Name)
	assert.Equal(t, newContent, inputChart.Templates[1].Data)
}

func TestChangeValues(t *testing.T) {

	meta := &chart.Metadata{
		Name: "test.chart",
	}
	originalContent := make(map[string]interface{})
	originalContent["image"] = "test:0.2"
	newContent := make(map[string]interface{})
	newContent["image"] = "test:latest"

	inputChart := chart.Chart{Metadata: meta}
	inputChart.Values = originalContent

	fileChanges := make(map[string]string)
	fileChanges["values.yaml"] = "image: test:latest"

	applyFileChanges(fileChanges, &inputChart)
	assert.Equal(t, newContent, inputChart.Values)
}
