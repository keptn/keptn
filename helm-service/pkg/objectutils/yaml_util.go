package objectutils

import "encoding/json"

// AppendAsYaml appends the element as yaml
func AppendAsYaml(content []byte, element interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(element)
	if err != nil {
		return nil, err
	}
	yamlData, err := ToYAML(jsonData)
	if err != nil {
		return nil, err
	}
	content = append(content, []byte("---\n")...)
	return append(content, yamlData...), nil
}

// AppendJSONStringAsYaml appends the JSON string as yaml
func AppendJSONStringAsYaml(content []byte, jsonString string) ([]byte, error) {

	yamlData, err := ToYAML([]byte(jsonString))
	if err != nil {
		return nil, err
	}
	content = append(content, []byte("---\n")...)
	return append(content, yamlData...), nil
}
