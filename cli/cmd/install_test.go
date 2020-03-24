package cmd

import (
	"reflect"
	"testing"
)

var iskeptnVersions = []struct {
	platform             string
	expectedPlatformType reflect.Type
	expectedErr          error
}{
	{"GKE", reflect.TypeOf(newGKEPlatform()), nil},
	{"AKS", reflect.TypeOf(newAKSPlatform()), nil},
	{"EKS", reflect.TypeOf(newEKSPlatform()), nil},
	{"PKS", reflect.TypeOf(newPKSPlatform()), nil},
	{"OpenShift", reflect.TypeOf(newOpenShiftPlatform()), nil},
	{"kubernetes", reflect.TypeOf(newKubernetesPlatform()), nil},
}

func TestSetPlatform(t *testing.T) {
	for _, tt := range iskeptnVersions {
		t.Run(tt.platform, func(t *testing.T) {
			installParams = &installCmdParams{PlatformIdentifier: &tt.platform}
			err := setPlatform()
			if err != tt.expectedErr {
				t.Errorf("got %t, want %t", err, tt.expectedErr)
			}
			if reflect.TypeOf(p) != tt.expectedPlatformType {
				t.Error("wrong type")
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func TestPrepareInstallerManifest(t *testing.T) {

	installParams = &installCmdParams{
		Image:              "keptn/installer",
		Tag:                "0.6.1",
		PlatformIdentifier: strPtr("gke"),
		Gateway:            LoadBalancer,
		UseCase:            AllUseCases,
		IstioInstallOption: StopIfInstalled,
	}

	res := prepareInstallerManifest()
	expected := `---
apiVersion: batch/v1
kind: Job
metadata:
  name: installer
  namespace: default
spec:
  backoffLimit: 0
  template:
    metadata:
      labels:
        app: installer
    spec:
      volumes:
      - name: kubectl
        emptyDir: {}
      containers:
      - name: keptn-installer
        image: keptn/installer:0.6.1
        env:
        - name: PLATFORM
          value: gke
        - name: GATEWAY_TYPE
          value: LoadBalancer
        - name: INGRESS
          value: istio
        - name: USE_CASE
          value: all
        - name: ISTIO_INSTALL_OPTION
          value: StopIfInstalled
      restartPolicy: Never
`
	if res != expected {
		t.Error("installation manifest does not match")
	}
}
