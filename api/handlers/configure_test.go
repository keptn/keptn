package handlers

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_getBridgeCredentials(t *testing.T) {
	type args struct {
		user     string
		password string
	}
	tests := []struct {
		name string
		args args
		want *corev1.Secret
	}{
		{
			name: "get bridge secret",
			args: args{
				user:     "user",
				password: "password",
			},
			want: &corev1.Secret{
				TypeMeta: v1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "apps/v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      "bridge-credentials",
					Namespace: "keptn",
				},
				Data: map[string][]byte{
					"BASIC_AUTH_USERNAME": []byte("user"),
					"BASIC_AUTH_PASSWORD": []byte("password"),
				},
				Type: "Opaque",
			},
		},
	}

	namespace = "keptn"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBridgeCredentials(tt.args.user, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBridgeCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}
