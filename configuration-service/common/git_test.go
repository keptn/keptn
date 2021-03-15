package common

import (
	"errors"
	"fmt"
	"github.com/bmizerany/assert"
	common_mock "github.com/keptn/keptn/configuration-service/common/fake"
	"github.com/keptn/keptn/configuration-service/common_models"
	"github.com/keptn/keptn/configuration-service/models"
	"strings"
	"testing"
)

func Test_obfuscateErrorMessage(t *testing.T) {
	type args struct {
		err         error
		credentials *common_models.GitCredentials
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantErrorMessage string
	}{
		{
			name: "remove credentials",
			args: args{
				err: errors.New("error message containing token: token"),
				credentials: &common_models.GitCredentials{
					User:      "",
					Token:     "token",
					RemoteURI: "",
				},
			},
			wantErr:          true,
			wantErrorMessage: "error message containing ********: ********",
		},
		{
			name: "remove credentials: empty token",
			args: args{
				err: errors.New("error message containing no token"),
				credentials: &common_models.GitCredentials{
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

func Test_addRepoURIToMetadata(t *testing.T) {
	type args struct {
		credentials *common_models.GitCredentials
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
				credentials: &common_models.GitCredentials{
					User:      "user",
					Token:     "secret-token",
					RemoteURI: "https://user:secret-token@my-url.test",
				},
				err: nil,
				resource: &models.Resource{
					Metadata: &models.Version{
						Branch: "master",
					},
					ResourceContent: "123",
					ResourceURI:     stringp("test.txt"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addRepoURIToMetadata(tt.args.credentials, tt.args.resource.Metadata)
			if strings.Contains(tt.args.resource.Metadata.UpstreamURL, tt.args.credentials.Token) {
				t.Errorf("Resource URI contains secret token")
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

func stringp(s string) *string {
	return &s
}

func TestGit_GetDefaultBranch(t *testing.T) {
	type fields struct {
		Executor         CommandExecutor
		CredentialReader CredentialReader
	}
	type args struct {
		project string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Get default branch - master branch available",
			fields: fields{
				Executor: &common_mock.CommandExecutorMock{
					ExecuteCommandFunc: func(command string, args []string, directory string) (string, error) {
						return `* remote origin
						  Fetch URL: https://my-repo.git
						  Push  URL: https://my-repo.git
						  HEAD branch: master
						  Remote branch:
							release-0.8.0 tracked
						  Local branch configured for 'git pull':
							release-0.8.0 merges with remote release-0.8.0
						  Local ref configured for 'git push':
							release-0.8.0 pushes to release-0.8.0 (up to date)`, nil
					},
				},
				CredentialReader: getDummyCredentialReader(),
			},
			args: args{
				project: "my-project",
			},
			want:    "master",
			wantErr: false,
		},
		{
			name: "Get default branch - no branch available at origin",
			fields: fields{
				Executor: &common_mock.CommandExecutorMock{
					ExecuteCommandFunc: func(command string, args []string, directory string) (string, error) {
						return `* remote origin
						  Fetch URL: https://my-repo.git
						  Push  URL: https://my-repo.git
						  HEAD branch: (unknown)`, nil
					},
				},
				CredentialReader: getDummyCredentialReader(),
			},
			args: args{
				project: "my-project",
			},
			want:    "master",
			wantErr: false,
		}, {
			name: "Get default branch - ambiguous HEAD",
			fields: fields{
				Executor: &common_mock.CommandExecutorMock{
					ExecuteCommandFunc: func(command string, args []string, directory string) (string, error) {
						if args[0] == "remote" {
							return `* remote origin
							  Fetch URL: https://my-repo.git
							  Push  URL: https://my-repo.git
							  HEAD branch (remote HEAD is ambiguous, may be one of the following):
								entry
								master
							  Remote branch:
								release-0.8.0 tracked
							  Local branch configured for 'git pull':
								release-0.8.0 merges with remote release-0.8.0
							  Local ref configured for 'git push':
								release-0.8.0 pushes to release-0.8.0 (up to date)`, nil
						} else if args[0] == "for-each-ref" {
							return fmt.Sprintf("dev\nmaster"), nil
						}
						return "", errors.New("unexpected command")
					},
				},
				CredentialReader: getDummyCredentialReader(),
			},
			args: args{
				project: "my-project",
			},
			want:    "master",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				Executor:         tt.fields.Executor,
				CredentialReader: tt.fields.CredentialReader,
			}
			got, err := g.GetDefaultBranch(tt.args.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDefaultBranch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getDummyCredentialReader() *common_mock.CredentialReaderMock {
	return &common_mock.CredentialReaderMock{
		GetCredentialsFunc: func(project string) (*common_models.GitCredentials, error) {
			return &common_models.GitCredentials{
				User:      "my-user",
				Token:     "token",
				RemoteURI: "https://my-repo.git",
			}, nil
		},
	}
}

func TestGit_setUpstreamsAndPush(t *testing.T) {
	type fields struct {
		Executor         *common_mock.CommandExecutorMock
		CredentialReader CredentialReader
	}
	type args struct {
		project     string
		credentials *common_models.GitCredentials
		repoURI     string
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantErr          bool
		expectedCommands []struct {
			Command   string
			Args      []string
			Directory string
		}
	}{
		{
			name: "push to upstream",
			fields: fields{
				Executor: &common_mock.CommandExecutorMock{ExecuteCommandFunc: func(command string, args []string, directory string) (string, error) {
					if args[0] == "for-each-ref" {
						return "master", nil
					} else if args[0] == "remote" {
						return `* remote origin
						  Fetch URL: https://my-repo.git
						  Push  URL: https://my-repo.git
						  HEAD branch: master
						  Remote branch:
							release-0.8.0 tracked
						  Local branch configured for 'git pull':
							release-0.8.0 merges with remote release-0.8.0
						  Local ref configured for 'git push':
							release-0.8.0 pushes to release-0.8.0 (up to date)`, nil
					} else if args[0] == "checkout" {
						return "", nil
					} else if args[0] == "push" {
						return "", nil
					}
					return "", errors.New("unexpected command")
				}},
				CredentialReader: getDummyCredentialReader(),
			},
			args: args{
				project: "my-project",
				repoURI: "https://my-repo.git",
			},
			wantErr: false,
			expectedCommands: []struct {
				Command   string
				Args      []string
				Directory string
			}{
				{
					Command:   "git",
					Args:      []string{"for-each-ref", "--format=%(refname:short)", "refs/heads/*"},
					Directory: "./debug/config/my-project",
				},
				{
					Command:   "git",
					Args:      []string{"remote", "show", "origin"},
					Directory: "./debug/config/my-project",
				},
				{
					Command:   "git",
					Args:      []string{"checkout", "master"},
					Directory: "./debug/config/my-project",
				},
				{
					Command:   "git",
					Args:      []string{"push", "--set-upstream", "https://my-repo.git", "master"},
					Directory: "./debug/config/my-project",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				Executor:         tt.fields.Executor,
				CredentialReader: tt.fields.CredentialReader,
			}
			if err := g.setUpstreamsAndPush(tt.args.project, tt.args.credentials, tt.args.repoURI); (err != nil) != tt.wantErr {
				t.Errorf("setUpstreamsAndPush() error = %v, wantErr %v", err, tt.wantErr)
			}
			executedCommands := tt.fields.Executor.ExecuteCommandCalls()

			assert.Equal(t, tt.expectedCommands, executedCommands)
		})
	}
}
