package common

import (
	"github.com/go-git/go-git/v5"
	"github.com/keptn/keptn/resource-service/common_models"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestK8sCredentialReader_GetCredentials(t *testing.T) {
	type args struct {
		project string
	}
	tests := []struct {
		name    string
		args    args
		want    *common_models.GitCredentials
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8 := K8sCredentialReader{}
			got, err := k8.GetCredentials(tt.args.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCredentials() got = %v, exists %v", got, tt.want)
			}
		})
	}
}

func Test_ensureRemoteMatchesCredentials(t *testing.T) {
	type args struct {
		repo       *git.Repository
		gitContext common_models.GitContext
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ensureRemoteMatchesCredentials(tt.args.repo, tt.args.gitContext); (err != nil) != tt.wantErr {
				t.Errorf("ensureRemoteMatchesCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getK8sClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *kubernetes.Clientset
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getK8sClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("getK8sClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getK8sClient() got = %v, exists %v", got, tt.want)
			}
		})
	}
}
