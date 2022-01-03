package go_tests

import (
	"testing"
)

func Test_K3D(t *testing.T) {
	// Common Tests
	t.Run("TestResourceServiceBasic", TestResourceServiceBasic)

	// Platform-specific Tests
	t.Run("TestAirgappedImagesAreSetCorrectly", TestAirgappedImagesAreSetCorrectly)
}
