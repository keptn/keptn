package models

import "testing"

func Test_validateResourceURI(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid resource name",
			args: args{
				uri: "helm/chart.tgz",
			},
			wantErr: false,
		},
		{
			name: "invalid resource name",
			args: args{
				uri: "../chart.tgz",
			},
			wantErr: true,
		},
		{
			name: "invalid resource name",
			args: args{
				uri: "~/chart.tgz",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateResourceURI(tt.args.uri); (err != nil) != tt.wantErr {
				t.Errorf("validateResourceURI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateResourcesParams_Validate(t *testing.T) {
	type fields struct {
		Project                Project
		Stage                  *Stage
		Service                *Service
		UpdateResourcesPayload UpdateResourcesPayload
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   &Stage{StageName: "my-stage"},
				Service: &Service{ServiceName: "my-service"},
				UpdateResourcesPayload: UpdateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "aGVsbG8K",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - not base64 encoded",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   &Stage{StageName: "my-stage"},
				Service: &Service{ServiceName: "my-service"},
				UpdateResourcesPayload: UpdateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "hello",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid project name",
			fields: fields{
				Project: Project{ProjectName: "my project"},
				Stage:   &Stage{StageName: "my-stage"},
				Service: &Service{ServiceName: "my-service"},
				UpdateResourcesPayload: UpdateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "aGVsbG8K",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid stage name",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   &Stage{StageName: "my stage"},
				Service: &Service{ServiceName: "my-service"},
				UpdateResourcesPayload: UpdateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "aGVsbG8K",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid service name",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   &Stage{StageName: "my-stage"},
				Service: &Service{ServiceName: "my service"},
				UpdateResourcesPayload: UpdateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "aGVsbG8K",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := UpdateResourcesParams{
				Project:                tt.fields.Project,
				Stage:                  tt.fields.Stage,
				Service:                tt.fields.Service,
				UpdateResourcesPayload: tt.fields.UpdateResourcesPayload,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateResourcesParams_Validate(t *testing.T) {
	type fields struct {
		Project                Project
		Stage                  *Stage
		Service                *Service
		CreateResourcesPayload CreateResourcesPayload
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   &Stage{StageName: "my-stage"},
				Service: &Service{ServiceName: "my-service"},
				CreateResourcesPayload: CreateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "aGVsbG8K",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - not base64 encoded",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   &Stage{StageName: "my-stage"},
				Service: &Service{ServiceName: "my-service"},
				CreateResourcesPayload: CreateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "hello",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid project name",
			fields: fields{
				Project: Project{ProjectName: "my project"},
				Stage:   &Stage{StageName: "my-stage"},
				Service: &Service{ServiceName: "my-service"},
				CreateResourcesPayload: CreateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "aGVsbG8K",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid stage name",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   &Stage{StageName: "my stage"},
				Service: &Service{ServiceName: "my-service"},
				CreateResourcesPayload: CreateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "aGVsbG8K",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid service name",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   &Stage{StageName: "my-stage"},
				Service: &Service{ServiceName: "my service"},
				CreateResourcesPayload: CreateResourcesPayload{
					Resources: []Resource{
						{
							ResourceContent: "aGVsbG8K",
							ResourceURI:     "resource.txt",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := CreateResourcesParams{
				Project:                tt.fields.Project,
				Stage:                  tt.fields.Stage,
				Service:                tt.fields.Service,
				CreateResourcesPayload: tt.fields.CreateResourcesPayload,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateResourceParams_Validate(t *testing.T) {
	type fields struct {
		Project               Project
		Stage                 *Stage
		Service               *Service
		ResourceURI           string
		UpdateResourcePayload UpdateResourcePayload
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Project:     Project{ProjectName: "my-project"},
				Stage:       &Stage{StageName: "my-stage"},
				Service:     &Service{ServiceName: "my-service"},
				ResourceURI: "my-resource.txt",
				UpdateResourcePayload: UpdateResourcePayload{
					ResourceContent: "aGVsbG8K",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := UpdateResourceParams{
				Project:               tt.fields.Project,
				Stage:                 tt.fields.Stage,
				Service:               tt.fields.Service,
				ResourceURI:           tt.fields.ResourceURI,
				UpdateResourcePayload: tt.fields.UpdateResourcePayload,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
