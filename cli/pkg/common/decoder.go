package common

import "github.com/mitchellh/mapstructure"

// DecodeKeptnEventData decodes the Data field of a Keptn Event to a given type
func DecodeKeptnEventData(in, out interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Squash: true,
		Result: out,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(in)
}
