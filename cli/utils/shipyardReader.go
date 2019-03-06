package utils

import (
	"log"

	"gopkg.in/yaml.v2"
)

type Stage map[string]string

func UnmarshalStages(data []byte) []Stage {

	stages := map[string][]Stage{}

	err := yaml.Unmarshal([]byte(data), &stages)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return stages["stages"]
}
