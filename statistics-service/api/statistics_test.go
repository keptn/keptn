package api

import (
	"github.com/keptn/keptn/statistics-service/controller"
	"github.com/keptn/keptn/statistics-service/db"
	"github.com/keptn/keptn/statistics-service/operations"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_validateQueryTimestamps(t *testing.T) {
	type args struct {
		params *operations.GetStatisticsParams
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "from < to",
			args: args{
				params: &operations.GetStatisticsParams{
					From: time.Now(),
					To:   time.Now().Add(1 * time.Second),
				},
			},
			want: true,
		},
		{
			name: "from < to",
			args: args{
				params: &operations.GetStatisticsParams{
					From: time.Now(),
					To:   time.Now().Add(-1 * time.Second),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateQueryTimestamps(tt.args.params); got != tt.want {
				t.Errorf("validateQueryTimestamps() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockStatisticsInterface struct {
	CutoffTime time.Time
	Statistics *operations.Statistics
	Repo       db.StatisticsRepo
}

func (m *MockStatisticsInterface) GetCutoffTime() time.Time {
	return m.CutoffTime
}

func (m *MockStatisticsInterface) GetStatistics() operations.Statistics {
	return *m.Statistics
}

func (m *MockStatisticsInterface) AddEvent(event operations.Event) {
	return
}

func (m *MockStatisticsInterface) GetRepo() db.StatisticsRepo {
	return m.Repo
}

// MockStatisticsRepo godoc
type MockStatisticsRepo struct {
	// GetStatisticsFunc godoc
	GetStatisticsFunc func(from, to time.Time) ([]operations.Statistics, error)
	// StoreStatisticsFunc godoc
	StoreStatisticsFunc func(statistics operations.Statistics) error
	// DeleteStatisticsFunc godoc
	DeleteStatisticsFunc func(from, to time.Time) error
}

// GetStatistics godoc
func (m *MockStatisticsRepo) GetStatistics(from, to time.Time) ([]operations.Statistics, error) {
	return m.GetStatisticsFunc(from, to)
}

// StoreStatistics godoc
func (m *MockStatisticsRepo) StoreStatistics(statistics operations.Statistics) error {
	return m.StoreStatisticsFunc(statistics)
}

// DeleteStatistics godoc
func (m *MockStatisticsRepo) DeleteStatistics(from, to time.Time) error {
	return m.DeleteStatisticsFunc(from, to)
}

func Test_getStatistics(t *testing.T) {
	type args struct {
		params     *operations.GetStatisticsParams
		statistics controller.StatisticsInterface
	}
	tests := []struct {
		name    string
		args    args
		want    operations.GetStatisticsResponse
		wantErr bool
	}{
		{
			name: "get in-memory bucket",
			args: args{
				params: &operations.GetStatisticsParams{
					From: time.Now(),
					To:   time.Now().Add(5 * time.Minute),
				},
				statistics: &MockStatisticsInterface{
					CutoffTime: time.Now().Add(-1 * time.Minute),
					Statistics: &operations.Statistics{
						From: time.Time{},
						To:   time.Time{},
						Projects: map[string]*operations.Project{
							"my-project": {
								Name:     "my-project",
								Services: map[string]*operations.Service{},
							},
						},
					},
					Repo: nil,
				},
			},
			want: operations.GetStatisticsResponse{
				From: time.Time{},
				To:   time.Time{},
				Projects: []operations.GetStatisticsResponseProject{
					{
						Name:     "my-project",
						Services: []operations.GetStatisticsResponseService{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "get bucket from db",
			args: args{
				params: &operations.GetStatisticsParams{
					From: time.Now().Round(time.Minute),
					To:   time.Now().Add(5 * time.Minute).Round(time.Minute),
				},
				statistics: &MockStatisticsInterface{
					CutoffTime: time.Now().Add(20 * time.Minute),
					Statistics: nil,
					Repo: &MockStatisticsRepo{
						GetStatisticsFunc: func(from, to time.Time) ([]operations.Statistics, error) {
							return []operations.Statistics{
								{
									From: time.Time{},
									To:   time.Time{},
									Projects: map[string]*operations.Project{
										"my-project": {
											Name:     "my-project",
											Services: map[string]*operations.Service{},
										},
									},
								},
								{
									From: time.Time{},
									To:   time.Time{},
									Projects: map[string]*operations.Project{
										"my-project-2": {
											Name:     "my-project-2",
											Services: nil,
										},
									},
								},
							}, nil
						},
					},
				},
			},
			want: operations.GetStatisticsResponse{
				From: time.Now().Round(time.Minute),
				To:   time.Now().Add(5 * time.Minute).Round(time.Minute),
				Projects: []operations.GetStatisticsResponseProject{
					{
						Name:     "my-project",
						Services: []operations.GetStatisticsResponseService{},
					},
					{
						Name:     "my-project-2",
						Services: []operations.GetStatisticsResponseService{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "get bucket from db and in-memory",
			args: args{
				params: &operations.GetStatisticsParams{
					From: time.Now().Round(time.Minute),
					To:   time.Now().Add(5 * time.Minute).Round(time.Minute),
				},
				statistics: &MockStatisticsInterface{
					CutoffTime: time.Now().Add(2 * time.Minute),
					Statistics: &operations.Statistics{
						From: time.Time{},
						To:   time.Time{},
						Projects: map[string]*operations.Project{
							"my-project-in-memory": {
								Name:     "my-project-in-memory",
								Services: map[string]*operations.Service{},
							},
						},
					},
					Repo: &MockStatisticsRepo{
						GetStatisticsFunc: func(from, to time.Time) ([]operations.Statistics, error) {
							return []operations.Statistics{
								{
									From: time.Time{},
									To:   time.Time{},
									Projects: map[string]*operations.Project{
										"my-project": {
											Name:     "my-project",
											Services: map[string]*operations.Service{},
										},
									},
								},
								{
									From: time.Time{},
									To:   time.Time{},
									Projects: map[string]*operations.Project{
										"my-project-2": {
											Name:     "my-project-2",
											Services: nil,
										},
									},
								},
							}, nil
						},
					},
				},
			},
			want: operations.GetStatisticsResponse{
				From: time.Now().Round(time.Minute),
				To:   time.Now().Add(5 * time.Minute).Round(time.Minute),
				Projects: []operations.GetStatisticsResponseProject{
					{
						Name:     "my-project",
						Services: []operations.GetStatisticsResponseService{},
					},
					{
						Name:     "my-project-2",
						Services: []operations.GetStatisticsResponseService{},
					},
					{
						Name:     "my-project-in-memory",
						Services: []operations.GetStatisticsResponseService{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStatistics(tt.args.params, tt.args.statistics)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStatistics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if match := assert.ElementsMatch(t, got.Projects, tt.want.Projects); !match {
				t.Errorf("getStatistics() got = %v, want %v", got, tt.want)
			}
		})
	}
}
