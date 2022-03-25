package validation

import (
	"github.com/keptn/keptn/shipyard-controller/config"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestProjectNameValidator(t *testing.T) {
	t.Run("Valid project name", func(t *testing.T) {
		v := NewProjectValidator(config.EnvConfig{ProjectNameMaxLength: 1})
		require.True(t, v.Validate(testFieldLevel{func() reflect.Value { return reflect.ValueOf("a") }}))
	})
	t.Run("Project name too long", func(t *testing.T) {
		v := NewProjectValidator(config.EnvConfig{ProjectNameMaxLength: 1})
		require.False(t, v.Validate(testFieldLevel{func() reflect.Value { return reflect.ValueOf("ab") }}))
	})

}

type testFieldLevel struct {
	FieldFn func() reflect.Value
}

func (t testFieldLevel) Top() reflect.Value {
	panic("implement me")
}

func (t testFieldLevel) Parent() reflect.Value {
	panic("implement me")
}

func (t testFieldLevel) Field() reflect.Value {
	return t.FieldFn()
}

func (t testFieldLevel) FieldName() string {
	panic("implement me")
}

func (t testFieldLevel) StructFieldName() string {
	panic("implement me")
}

func (t testFieldLevel) Param() string {
	panic("implement me")
}

func (t testFieldLevel) GetTag() string {
	panic("implement me")
}

func (t testFieldLevel) ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool) {
	panic("implement me")
}

func (t testFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	panic("implement me")
}

func (t testFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	panic("implement me")
}

func (t testFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	panic("implement me")
}

func (t testFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	panic("implement me")
}
