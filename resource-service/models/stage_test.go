package models

import "testing"

func TestDeleteStageParams_Validate(t *testing.T) {
	type fields struct {
		Project Project
		Stage   Stage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Project: Project{
					ProjectName: "my-project",
				},
				Stage: Stage{
					StageName: "my-stage",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid project",
			fields: fields{
				Project: Project{
					ProjectName: "my project",
				},
				Stage: Stage{
					StageName: "my-stage",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid stage",
			fields: fields{
				Project: Project{
					ProjectName: "my-project",
				},
				Stage: Stage{
					StageName: "my stage",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DeleteStageParams{
				Project: tt.fields.Project,
				Stage:   tt.fields.Stage,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateStageParams_Validate(t *testing.T) {
	type fields struct {
		Project            Project
		CreateStagePayload CreateStagePayload
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Project: Project{
					ProjectName: "my-project",
				},
				CreateStagePayload: CreateStagePayload{
					Stage{
						StageName: "my-stage",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid project",
			fields: fields{
				Project: Project{
					ProjectName: "my project",
				},
				CreateStagePayload: CreateStagePayload{
					Stage{
						StageName: "my-stage",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid stage",
			fields: fields{
				Project: Project{
					ProjectName: "my-project",
				},
				CreateStagePayload: CreateStagePayload{
					Stage{
						StageName: "my stage",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CreateStageParams{
				Project:            tt.fields.Project,
				CreateStagePayload: tt.fields.CreateStagePayload,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
