package models

import "testing"

func TestDeleteProjectPathParams_Validate(t *testing.T) {
	type fields struct {
		Project Project
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid name+",
			fields: fields{
				Project: Project{
					ProjectName: "my-project",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid name+",
			fields: fields{
				Project: Project{
					ProjectName: "my project",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := DeleteProjectPathParams{
				Project: tt.fields.Project,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateProjectParams_Validate(t *testing.T) {
	type fields struct {
		Project Project
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid name+",
			fields: fields{
				Project: Project{
					ProjectName: "my-project",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid name+",
			fields: fields{
				Project: Project{
					ProjectName: "my project",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := UpdateProjectParams{
				Project: tt.fields.Project,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateProjectParams_Validate(t *testing.T) {
	type fields struct {
		Project Project
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid name+",
			fields: fields{
				Project: Project{
					ProjectName: "my-project",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid name+",
			fields: fields{
				Project: Project{
					ProjectName: "my project",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := CreateProjectParams{
				Project: tt.fields.Project,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
