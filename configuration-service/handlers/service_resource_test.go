package handlers

import "testing"

func Test_isHelmChart(t *testing.T) {
	type args struct {
		resourcePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "helm chart file path",
			args: args{
				resourcePath: "helm/chart.tgz",
			},
			want: true,
		},
		{
			name: "not helm chart file path",
			args: args{
				resourcePath: "helm/chart.tgza",
			},
			want: false,
		},
		{
			name: "not helm chart file path",
			args: args{
				resourcePath: "helmy/chart.tgz",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHelmChart(tt.args.resourcePath); got != tt.want {
				t.Errorf("isHelmChart() = %v, want %v", got, tt.want)
			}
		})
	}
}
