package objectutils

import (
	"bytes"
	"encoding/json"

	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// ToJSON converts a yaml byte slice into a json byte slice
func ToJSON(yaml []byte) ([]byte, error) {
	var jsonData interface{}
	dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(yaml))
	dec.Decode(&jsonData)
	return json.Marshal(jsonData)
}

// ToYAML converts a json byte slice into a yaml byte slice
func ToYAML(json []byte) ([]byte, error) {
	return yaml.JSONToYAML(json)
}
