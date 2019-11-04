package utils

import (
	"gopkg.in/yaml.v2"
)

// UnmarshalString provides a YAML unmarsheling helper
func UnmarshalString(data string) (interface{}, error) {
	var body interface{}
	err := yaml.Unmarshal([]byte(data), &body)
	if err != nil {
		return nil, err
	}
	return Convert(body), nil
}

// Convert makes a type conversion of a yaml object
func Convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = Convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = Convert(v)
		}
	}
	return i
}
