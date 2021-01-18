package helm

import (
	"os"
	"testing"
)

func Test_getVirtualServicePublicHost(t *testing.T) {
	type args struct {
		svc       string
		project   string
		stageName string
	}
	tests := []struct {
		name             string
		hostnameTemplate string
		hostnameSuffix   string
		args             args
		want             string
		wantErr          bool
	}{
		{
			name: "get default hostname",
			args: args{
				svc:       "svc",
				project:   "prj",
				stageName: "stg",
			},
			want:    "svc.prj-stg.svc.cluster.local",
			wantErr: false,
		},
		{
			name:           "get hostname based on default template and custom INGRESS_HOSTNAME_SUFFIX",
			hostnameSuffix: "123.xip.io",
			args: args{
				svc:       "svc",
				project:   "prj",
				stageName: "stg",
			},
			want:    "svc.prj-stg.123.xip.io",
			wantErr: false,
		},
		{
			name:             "get hostname based on custom HOSTNAME_TEMPLATE and custom INGRESS_HOSTNAME_SUFFIX",
			hostnameTemplate: "${service}-${stage}-${project}.${INGRESS_HOSTNAME_SUFFIX}",
			hostnameSuffix:   "123.xip.io",
			args: args{
				svc:       "svc",
				project:   "prj",
				stageName: "stg",
			},
			want:    "",
			wantErr: true,
		},
		{
			name:             "get hostname based on custom HOSTNAME_TEMPLATE and custom INGRESS_HOSTNAME_SUFFIX",
			hostnameTemplate: "${INGRESS_PROTOCOL}://${service}-${stage}-${project}.${INGRESS_HOSTNAME_SUFFIX}",
			hostnameSuffix:   "123.xip.io",
			args: args{
				svc:       "svc",
				project:   "prj",
				stageName: "stg",
			},
			want:    "svc-stg-prj.123.xip.io",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOSTNAME_TEMPLATE", tt.hostnameTemplate)
			os.Setenv("INGRESS_HOSTNAME_SUFFIX", tt.hostnameSuffix)
			got, err := getVirtualServicePublicHost(tt.args.svc, tt.args.project, tt.args.stageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVirtualServicePublicHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getVirtualServicePublicHost() got = %v, want %v", got, tt.want)
			}
		})
	}
}
