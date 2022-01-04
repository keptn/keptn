package common

import "testing"

func TestGetProjectConfigPath(t *testing.T) {
	type args struct {
		project string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get project path",
			args: args{
				project: "my-project",
			},
			want: ConfigDir + "/my-project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjectConfigPath(tt.args.project); got != tt.want {
				t.Errorf("GetProjectConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProjectMetadataFilePath(t *testing.T) {
	type args struct {
		project string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get project path",
			args: args{
				project: "my-project",
			},
			want: ConfigDir + "/my-project/metadata.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjectMetadataFilePath(tt.args.project); got != tt.want {
				t.Errorf("GetProjectMetadataFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceConfigPath(t *testing.T) {
	type args struct {
		project string
		service string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get project path",
			args: args{
				project: "my-project",
				service: "my-service",
			},
			want: ConfigDir + "/my-project/my-service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceConfigPath(tt.args.project, tt.args.service); got != tt.want {
				t.Errorf("GetServiceConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
