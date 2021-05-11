package common

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetKeptnSpecVersion(t *testing.T) {
	specVersion := GetKeptnSpecVersion()
	assert.Equal(t, "", specVersion)

	os.Setenv(keptnSpecVersionEnvVar, "0.2.0")
	specVersion = GetKeptnSpecVersion()
	assert.Equal(t, "0.2.0", specVersion)
}

func TestValidateShipyardVersion(t *testing.T) {
	type args struct {
		shipyard *keptnv2.Shipyard
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid shipyard version",
			args: args{
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "0.2.0",
				},
			},
			wantErr: false,
		},
		{
			name: "valid shipyard version 2",
			args: args{
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "spec.keptn.sh/0.2.0",
				},
			},
			wantErr: false,
		},
		{
			name: "valid shipyard version 3",
			args: args{
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "0.2.2",
				},
			},
			wantErr: false,
		},
		{
			name: "valid shipyard version 4",
			args: args{
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "spec.keptn.sh/0.2.2",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid shipyard version",
			args: args{
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "spec.keptn.sh/0.1.0",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateShipyardVersion(tt.args.shipyard); (err != nil) != tt.wantErr {
				t.Errorf("ValidateShipyardVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
