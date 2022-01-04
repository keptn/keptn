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
				ResourceContext: ResourceContext{
					Project: tt.fields.Project,
					Stage:   tt.fields.Stage,
					Service: tt.fields.Service,
				},
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
				ResourceContext: ResourceContext{
					Project: tt.fields.Project,
					Stage:   tt.fields.Stage,
					Service: tt.fields.Service,
				},
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
		{
			name: "invalid",
			fields: fields{
				Project:     Project{ProjectName: "my project"},
				Stage:       &Stage{StageName: "my-stage"},
				Service:     &Service{ServiceName: "my-service"},
				ResourceURI: "my-resource.txt",
				UpdateResourcePayload: UpdateResourcePayload{
					ResourceContent: "aGVsbG8K",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid",
			fields: fields{
				Project:     Project{ProjectName: "my-project"},
				Stage:       &Stage{StageName: "my stage"},
				Service:     &Service{ServiceName: "my-service"},
				ResourceURI: "my-resource.txt",
				UpdateResourcePayload: UpdateResourcePayload{
					ResourceContent: "aGVsbG8K",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid",
			fields: fields{
				Project:     Project{ProjectName: "my-project"},
				Stage:       &Stage{StageName: "my-stage"},
				Service:     &Service{ServiceName: "my service"},
				ResourceURI: "my-resource.txt",
				UpdateResourcePayload: UpdateResourcePayload{
					ResourceContent: "aGVsbG8K",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid",
			fields: fields{
				Project:     Project{ProjectName: "my-project"},
				Stage:       &Stage{StageName: "my-stage"},
				Service:     &Service{ServiceName: "my-service"},
				ResourceURI: "../my-resource.txt",
				UpdateResourcePayload: UpdateResourcePayload{
					ResourceContent: "aGVsbG8K",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid",
			fields: fields{
				Project:     Project{ProjectName: "my-project"},
				Stage:       &Stage{StageName: "my-stage"},
				Service:     &Service{ServiceName: "my-service"},
				ResourceURI: "my-resource.txt",
				UpdateResourcePayload: UpdateResourcePayload{
					ResourceContent: "test123",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := UpdateResourceParams{
				ResourceContext: ResourceContext{
					Project: tt.fields.Project,
					Stage:   tt.fields.Stage,
					Service: tt.fields.Service,
				},
				ResourceURI:           tt.fields.ResourceURI,
				UpdateResourcePayload: tt.fields.UpdateResourcePayload,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteResourceParams_Validate(t *testing.T) {
	type fields struct {
		ResourceContext ResourceContext
		ResourceURI     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ResourceContext: ResourceContext{
					Project: Project{ProjectName: "my-project"},
					Stage:   &Stage{StageName: "my-stage"},
					Service: &Service{ServiceName: "my-service"},
				},
				ResourceURI: "my-resource.txt",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				ResourceContext: ResourceContext{
					Project: Project{ProjectName: "my project"},
					Stage:   &Stage{StageName: "my-stage"},
					Service: &Service{ServiceName: "my-service"},
				},
				ResourceURI: "my-resource.txt",
			},
			wantErr: true,
		},
		{
			name: "invalid",
			fields: fields{
				ResourceContext: ResourceContext{
					Project: Project{ProjectName: "my-project"},
					Stage:   &Stage{StageName: "my-stage"},
					Service: &Service{ServiceName: "my-service"},
				},
				ResourceURI: "../my-resource.txt",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := DeleteResourceParams{
				ResourceContext: tt.fields.ResourceContext,
				ResourceURI:     tt.fields.ResourceURI,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetResourceParams_Validate(t *testing.T) {
	type fields struct {
		ResourceContext  ResourceContext
		ResourceURI      string
		GetResourceQuery GetResourceQuery
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ResourceContext: ResourceContext{
					Project: Project{ProjectName: "my-project"},
					Stage:   &Stage{StageName: "my-stage"},
					Service: &Service{ServiceName: "my-service"},
				},
				ResourceURI: "my-resource.txt",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				ResourceContext: ResourceContext{
					Project: Project{ProjectName: "my project"},
					Stage:   &Stage{StageName: "my-stage"},
					Service: &Service{ServiceName: "my-service"},
				},
				ResourceURI: "my-resource.txt",
			},
			wantErr: true,
		},
		{
			name: "invalid",
			fields: fields{
				ResourceContext: ResourceContext{
					Project: Project{ProjectName: "my-project"},
					Stage:   &Stage{StageName: "my-stage"},
					Service: &Service{ServiceName: "my-service"},
				},
				ResourceURI: "../my-resource.txt",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := GetResourceParams{
				ResourceContext:  tt.fields.ResourceContext,
				ResourceURI:      tt.fields.ResourceURI,
				GetResourceQuery: tt.fields.GetResourceQuery,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetResourcesParams_Validate(t *testing.T) {
	type fields struct {
		ResourceContext   ResourceContext
		GetResourcesQuery GetResourcesQuery
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ResourceContext: ResourceContext{
					Project: Project{ProjectName: "my-project"},
					Stage:   &Stage{StageName: "my-stage"},
					Service: &Service{ServiceName: "my-service"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				ResourceContext: ResourceContext{
					Project: Project{ProjectName: "my project"},
					Stage:   &Stage{StageName: "my-stage"},
					Service: &Service{ServiceName: "my-service"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := GetResourcesParams{
				ResourceContext:   tt.fields.ResourceContext,
				GetResourcesQuery: tt.fields.GetResourcesQuery,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResource_Validate(t *testing.T) {
	type fields struct {
		ResourceContent ResourceContent
		ResourceURI     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ResourceContent: "aGVsbG8K",
				ResourceURI:     "resource.txt",
			},
			wantErr: false,
		},
		{
			name: "invalid uri",
			fields: fields{
				ResourceContent: "aGVsbG8K",
				ResourceURI:     "../resource.txt",
			},
			wantErr: true,
		},
		{
			name: "invalid content",
			fields: fields{
				ResourceContent: "123",
				ResourceURI:     "resource.txt",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Resource{
				ResourceContent: tt.fields.ResourceContent,
				ResourceURI:     tt.fields.ResourceURI,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
