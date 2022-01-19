package models

import "testing"

func TestDeleteServiceParams_Validate(t *testing.T) {
	type fields struct {
		Project Project
		Stage   Stage
		Service Service
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
				Stage:   Stage{StageName: "my-stage"},
				Service: Service{ServiceName: "my-service"},
			},
			wantErr: false,
		},
		{
			name: "invalid project",
			fields: fields{
				Project: Project{ProjectName: "my project"},
				Stage:   Stage{StageName: "my-stage"},
				Service: Service{ServiceName: "my-service"},
			},
			wantErr: true,
		},
		{
			name: "invalid stage",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   Stage{StageName: "my stage"},
				Service: Service{ServiceName: "my-service"},
			},
			wantErr: true,
		},
		{
			name: "invalid service",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   Stage{StageName: "my-stage"},
				Service: Service{ServiceName: "my service"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DeleteServiceParams{
				Project: tt.fields.Project,
				Stage:   tt.fields.Stage,
				Service: tt.fields.Service,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateServiceParams_Validate(t *testing.T) {
	type fields struct {
		Project              Project
		Stage                Stage
		CreateServicePayload CreateServicePayload
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
				Stage:   Stage{StageName: "my-stage"},
				CreateServicePayload: CreateServicePayload{
					Service: Service{
						ServiceName: "my-service",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid project",
			fields: fields{
				Project: Project{ProjectName: "my project"},
				Stage:   Stage{StageName: "my-stage"},
				CreateServicePayload: CreateServicePayload{
					Service: Service{
						ServiceName: "my-service",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid stage",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   Stage{StageName: "my stage"},
				CreateServicePayload: CreateServicePayload{
					Service: Service{
						ServiceName: "my-service",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid service",
			fields: fields{
				Project: Project{ProjectName: "my-project"},
				Stage:   Stage{StageName: "my-stage"},
				CreateServicePayload: CreateServicePayload{
					Service: Service{
						ServiceName: "my service",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CreateServiceParams{
				Project:              tt.fields.Project,
				Stage:                tt.fields.Stage,
				CreateServicePayload: tt.fields.CreateServicePayload,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
