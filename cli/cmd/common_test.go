package cmd

import "testing"

func Test_isStringFlagSet(t *testing.T) {
	type args struct {
		s *string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "flag is set and not empty - return true",
			args: args{
				s: stringp("foo"),
			},
			want: true,
		},
		{
			name: "flag is set but empty - return false",
			args: args{
				s: stringp(""),
			},
			want: false,
		},
		{
			name: "flag is nil - return false",
			args: args{
				s: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStringFlagSet(tt.args.s); got != tt.want {
				t.Errorf("isStringFlagSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isBoolFlagSet(t *testing.T) {
	var trueFlag, falseFlag bool
	trueFlag = true
	falseFlag = false
	type args struct {
		s *bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "flag is set and true - return true",
			args: args{
				s: &trueFlag,
			},
			want: true,
		},
		{
			name: "flag is set but false - return false",
			args: args{
				s: &falseFlag,
			},
			want: false,
		},
		{
			name: "flag is nil - return false",
			args: args{
				s: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBoolFlagSet(tt.args.s); got != tt.want {
				t.Errorf("isStringFlagSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_areStringFlagsSet(t *testing.T) {
	type args struct {
		el []*string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all flags set - return true",
			args: args{
				el: []*string{stringp("foo"), stringp("bar")},
			},
			want: true,
		},
		{
			name: "not all flags set - return false",
			args: args{
				el: []*string{stringp("foo"), stringp("")},
			},
			want: false,
		},
		{
			name: "not all flags set - return false",
			args: args{
				el: []*string{stringp("foo"), nil},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := areStringFlagsSet(tt.args.el...); got != tt.want {
				t.Errorf("areStringFlagsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_areBoolFlagsSet(t *testing.T) {
	var trueFlag, falseFlag bool
	trueFlag = true
	falseFlag = false
	type args struct {
		el []*bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all flags set - return true",
			args: args{
				el: []*bool{&trueFlag, &trueFlag},
			},
			want: true,
		},
		{
			name: "not all flags set - return false",
			args: args{
				el: []*bool{&trueFlag, &falseFlag},
			},
			want: false,
		},
		{
			name: "not all flags set - return false",
			args: args{
				el: []*bool{&trueFlag, nil},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := areBoolFlagsSet(tt.args.el...); got != tt.want {
				t.Errorf("areBoolFlagsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
