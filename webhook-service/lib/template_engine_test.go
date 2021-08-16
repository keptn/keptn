package lib_test

import (
	"github.com/keptn/keptn/webhook-service/lib"
	"testing"
)

func TestTemplateEngine_ParseTemplate(t1 *testing.T) {
	type args struct {
		data        interface{}
		templateStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "parse template",
			args: args{
				data: map[string]interface{}{
					"env": map[string]interface{}{
						"bar": "bar",
					},
				},
				templateStr: "foo {{.env.bar}}",
			},
			want:    "foo bar",
			wantErr: false,
		},
		{
			name: "wrong template syntax",
			args: args{
				data: map[string]interface{}{
					"env": map[string]interface{}{
						"bar": "bar",
					},
				},
				templateStr: "foo {{.Env.Bar}",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &lib.TemplateEngine{}
			got, err := t.ParseTemplate(tt.args.data, tt.args.templateStr)
			if (err != nil) != tt.wantErr {
				t1.Errorf("ParseTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("ParseTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
