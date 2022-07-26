package projectcreator

import (
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/internal/projectcreator/fake"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestProjectCreator_CreateProject(t *testing.T) {

	type fields struct {
		APIV1Interface   api.APIV1Interface
		ShipyardProvider ShipyardProvider
		FileSystem       fs.FS
	}
	type args struct {
		projectInfo ProjectInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Getting MetaData fails",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return nil, &apimodels.Error{}
					},
				},
			},
			args:    args{projectInfo: ProjectInfo{}},
			wantErr: true,
		},
		{
			name: "Failed to get shipyard file",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return nil, fmt.Errorf("fail") },
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:  "shipyard-file",
				GitUser:   "gitUser",
				RemoteURL: "remotURL",
			}},
			wantErr: true,
		},
		{
			name: "Access token or private key must be set",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:  "shipyard-file",
				GitUser:   "gitUser",
				RemoteURL: "remotURL",
			}},
			wantErr: true,
		},
		{
			name: "Access token and private key cannot be set together",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:      "shipyard-file",
				GitToken:      "git-token",
				GitPrivateKey: "private-key",
				GitUser:       "gitUser",
				RemoteURL:     "remotURL",
			}},
			wantErr: true,
		},
		{
			name: "Proxy cannot be set with SSH",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:    "shipyard-file",
				GitToken:    "git-token",
				GitUser:     "gitUser",
				RemoteURL:   "ssh://remote-url.com",
				GitProxyURL: "http://some-url.com",
			}},
			wantErr: true,
		},
		{
			name: "Proxy cannot be set with SSH",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:    "shipyard-file",
				GitToken:    "git-token",
				GitUser:     "gitUser",
				RemoteURL:   "ssh://remote-url.com",
				GitProxyURL: "http://some-url.com",
			}},
			wantErr: true,
		},
		{
			name: "Proxy cannot be set without scheme",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:    "shipyard-file",
				GitToken:    "git-token",
				GitUser:     "gitUser",
				RemoteURL:   "http://remote-url.com",
				GitProxyURL: "http://some-url.com",
			}},
			wantErr: true,
		},
		{
			name: "Create Project with GIT information - Creating project fails",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					CreateProjectFunc: func(project apimodels.CreateProject) (string, *apimodels.Error) {
						return "", &apimodels.Error{}
					},
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:         "shipyard-file",
				GitToken:         "git-token",
				GitUser:          "gitUser",
				RemoteURL:        "http://remote-url.com",
				GitProxyURL:      "http://some-url.com",
				GitProxyScheme:   "scheme",
				GitProxyUser:     "proxy-user",
				GitProxyPassword: "proxy-pwd",
				InsecureSkipTLS:  false,
			}},
			wantErr: true,
		},
		{
			name: "Create Project with SSH private key",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
					CreateProjectFunc: func(project apimodels.CreateProject) (string, *apimodels.Error) {
						return "project", nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
				FileSystem: fstest.MapFS{
					"git-private-key-path": {Data: []byte("private-key")},
				},
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:          "shipyard-file",
				GitUser:           "gitUser",
				GitPrivateKey:     "git-private-key-path",
				GitPrivateKeyPass: "git-private-key-pass",
				RemoteURL:         "ssh://remote-url.com",
			}},
			wantErr: false,
		},
		{
			name: "Create Project with SSH private key - unable to read private key file",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
					CreateProjectFunc: func(project apimodels.CreateProject) (string, *apimodels.Error) {
						return "project", nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
				FileSystem:       fstest.MapFS{},
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:          "shipyard-file",
				GitUser:           "gitUser",
				GitPrivateKey:     "git-private-key-path",
				GitPrivateKeyPass: "git-private-key-pass",
				RemoteURL:         "ssh://remote-url.com",
			}},
			wantErr: true,
		},
		{
			name: "Create Project with certificate",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
					CreateProjectFunc: func(project apimodels.CreateProject) (string, *apimodels.Error) {
						return "project", nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
				FileSystem: fstest.MapFS{
					"pem-certificate-path": {Data: []byte("certificate")},
				},
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:          "shipyard-file",
				GitUser:           "gitUser",
				GitToken:          "token",
				GitPemCertificate: "pem-certificate-path",
				InsecureSkipTLS:   true,
				RemoteURL:         "http://remote-url.com",
			}},
			wantErr: false,
		},
		{
			name: "Create Project with certificate - reading certificate fails",
			fields: fields{
				APIV1Interface: &fake.APIV1InterfaceMock{
					GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
						return &apimodels.Metadata{
							Automaticprovisioning: boolP(false),
						}, nil
					},
					CreateProjectFunc: func(project apimodels.CreateProject) (string, *apimodels.Error) {
						return "project", nil
					},
				},
				ShipyardProvider: func(s string) ([]byte, error) { return []byte{}, nil },
				FileSystem:       fstest.MapFS{},
			},
			args: args{projectInfo: ProjectInfo{
				Shipyard:          "shipyard-file",
				GitUser:           "gitUser",
				GitToken:          "token",
				GitPemCertificate: "pem-certificate-path",
				InsecureSkipTLS:   true,
				RemoteURL:         "http://remote-url.com",
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.fields.APIV1Interface, tt.fields.ShipyardProvider, tt.fields.FileSystem)

			if err := p.CreateProject(tt.args.projectInfo); (err != nil) != tt.wantErr {
				t.Errorf("CreateProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func boolP(b bool) *bool {
	return &b
}

func stringP(s string) *string {
	return &s
}
