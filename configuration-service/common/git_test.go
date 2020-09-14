package common

import (
	"errors"
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

func Test_getRepoURI(t *testing.T) {
	type args struct {
		uri   string
		user  string
		token string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get url with https:// ",
			args: args{
				uri:   "https://my-repo.git",
				user:  "user",
				token: "token",
			},
			want: "https://user:token@my-repo.git",
		},
		{
			name: "get url with https:// where user is already included in the url",
			args: args{
				uri:   "https://user@my-repo.git",
				user:  "user",
				token: "token",
			},
			want: "https://user:token@my-repo.git",
		},
		{
			name: "get url with http:// ",
			args: args{
				uri:   "http://my-repo.git",
				user:  "user",
				token: "token",
			},
			want: "http://user:token@my-repo.git",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRepoURI(tt.args.uri, tt.args.user, tt.args.token); got != tt.want {
				t.Errorf("getRepoURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
