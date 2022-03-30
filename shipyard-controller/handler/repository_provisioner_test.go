package handler

import (
	"testing"
)

func TestProvideRepository(t *testing.T) {
	tests := []struct {
		name                  string
		RepositoryProvisioner *RepositoryProvisioner
		wantErr               bool
	}{
		{
			name:                  "basic test",
			RepositoryProvisioner: NewRepositoryProvisioner("som-url"),
			wantErr:               true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.RepositoryProvisioner.ProvideRepository("testing-project")
			if (err != nil) != tt.wantErr {
				t.Errorf("ProvideRepository() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteRepository(t *testing.T) {
	tests := []struct {
		name                  string
		RepositoryProvisioner *RepositoryProvisioner
		wantErr               bool
	}{
		{
			name:                  "basic test",
			RepositoryProvisioner: NewRepositoryProvisioner("som-url"),
			wantErr:               true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.RepositoryProvisioner.DeleteRepository("testing-project", "testing-namespace")
			if (err != nil) != tt.wantErr {
				t.Errorf("ProvideRepository() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
