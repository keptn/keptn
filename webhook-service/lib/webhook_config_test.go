package lib

import (
	"reflect"
	"testing"
)

func TestDecodeWebHookConfigYAML(t *testing.T) {
	type args struct {
		webhookConfigYaml []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *WebHookConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeWebHookConfigYAML(tt.args.webhookConfigYaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeWebHookConfigYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeWebHookConfigYAML() got = %v, want %v", got, tt.want)
			}
		})
	}
}
