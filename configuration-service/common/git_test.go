package common

import (
	"errors"
	"github.com/keptn/keptn/configuration-service/models"
	"strings"
	"testing"
)

func Test_obfuscateErrorMessage(t *testing.T) {
	type args struct {
		err         error
		credentials *GitCredentials
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantErrorMessage string
	}{
		{
			name: "remnove credentials",
			args: args{
				err: errors.New("error message containing token: token"),
				credentials: &GitCredentials{
					User:      "",
					Token:     "token",
					RemoteURI: "",
				},
			},
			wantErr:          true,
			wantErrorMessage: "error message containing ********: ********",
		},
		{
			name: "remnove credentials: empty token",
			args: args{
				err: errors.New("error message containing no token"),
				credentials: &GitCredentials{
					User:      "",
					Token:     "",
					RemoteURI: "",
				},
			},
			wantErr:          true,
			wantErrorMessage: "error message containing no token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := obfuscateErrorMessage(tt.args.err, tt.args.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("obfuscateErrorMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err.Error() != tt.wantErrorMessage {
				t.Errorf("obfuscateErrorMessage() error = %s, wantErrorMessage %s", err.Error(), tt.wantErrorMessage)
			}
		})
	}
}

func Test_addRepoURIToResource(t *testing.T) {
	type args struct {
		credentials *GitCredentials
		err         error
		resource    *models.Resource
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Token must not be included",
			args: args{
				credentials: &GitCredentials{
					User:      "user",
					Token:     "secret-token",
					RemoteURI: "https://user:secret-token@my-url.test",
				},
				err: nil,
				resource: &models.Resource{
					Branch:          "master",
					ResourceContent: "123",
					ResourceURI:     stringp("test.txt"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addRepoURIToResource(tt.args.credentials, tt.args.resource)
			if strings.Contains(tt.args.resource.UpstreamURL, tt.args.credentials.Token) {
				t.Errorf("Resource URI contains secret token")
			}
		})
	}
}
