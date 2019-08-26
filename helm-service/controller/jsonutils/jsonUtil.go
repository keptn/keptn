package jsonutils

import (
	"bytes"
	"encoding/json"

	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

// ToJSON converts a yaml byte slice into a json byte slice
func ToJSON(yaml []byte) ([]byte, error) {
	var jsonData interface{}
	dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(yaml))
	dec.Decode(&jsonData)
	return json.Marshal(jsonData)
}
