package common

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_EnvBasedStringSupplier(t *testing.T) {
	os.Setenv("KEPTN_TEST_ENV_VAR", "KEPTN_TEST_ENV_VAR_VAL")
	val := EnvBasedStringSupplier("KEPTN_TEST_ENV_VAR", "")()
	assert.Equal(t, "KEPTN_TEST_ENV_VAR_VAL", val)

	val = EnvBasedStringSupplier("THIS_ENV_VAR_IS_NOT_PRESENT", "DEFAULT_VAL")()
	assert.Equal(t, "DEFAULT_VAL", val)
}
