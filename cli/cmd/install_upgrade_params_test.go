package cmd

import "testing"

func Test_getServiceChartURLFromKeptnChartURL(t *testing.T) {
	type args struct {
		input       string
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get helm-service chart URL",
			args: args{
				input:       "https://keptn.sh/keptn-0.8.1.tgz",
				serviceName: "helm-service",
			},
			want: "https://keptn.sh/helm-service-0.8.1.tgz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getServiceChartURLFromKeptnChartURL(tt.args.input, tt.args.serviceName); got != tt.want {
				t.Errorf("getServiceChartURLFromKeptnChartURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
