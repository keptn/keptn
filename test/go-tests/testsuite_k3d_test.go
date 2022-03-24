package go_tests

import (
	"testing"
)

func Test_K3D(t *testing.T) {
	// Platform-specific Tests
	t.Run("TestAirgappedImagesAreSetCorrectly", TestAirgappedImagesAreSetCorrectly)
}
