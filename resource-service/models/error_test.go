package models

import "testing"

func TestError_Error(t *testing.T) {
	type fields struct {
		Code    int64
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get error message",
			fields: fields{
				Code:    0,
				Message: "error msg",
			},
			want: "error msg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
			}
			if got := m.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
