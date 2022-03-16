package go_tests

import "testing"

// Test_GKE_MR contains tests that are run against a keptn installation
// that has multiple replicas of its components in place
func Test_GKE_MR(t *testing.T) {
	t.Run("Test_ResourceService_MR", Test_ResourceServiceBasic)
}
