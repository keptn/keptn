package go_tests

import (
	"testing"
)

func Test_K3D(t *testing.T) {
	// Common Tests
	t.Run("Test_ResourceServiceBasic", Test_ResourceServiceBasic)

	// Platform-specific Tests
	t.Run("TestAirgappedImagesAreSetCorrectly", TestAirgappedImagesAreSetCorrectly)
}
