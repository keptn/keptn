package go_tests

import (
	"testing"
)

func Test_K3S(t *testing.T) {
	// Common Tests
	t.Run("TestCommon", TestCommon)

	// Platform-specific Tests
	t.Run("Test_QualityGates", Test_QualityGates)
	t.Run("Test_QualityGates_BackwardsCompatibility", Test_QualityGates_BackwardsCompatibility)
	t.Run("Test_UniformRegistration_TestAPI", Test_UniformRegistration_TestAPI)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegration", Test_UniformRegistration_RegistrationOfKeptnIntegration)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods", Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods)
	t.Run("Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane", Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane)
}
