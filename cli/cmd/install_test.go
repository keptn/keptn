package cmd

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
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
			*installParams.PlatformIdentifier = tt.platform
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

func TestPrepareInstallerManifest(t *testing.T) {

	*installParams.InstallerImage = "keptn/installer:0.6.1"
	*installParams.PlatformIdentifier = "gke"
	installParams.Gateway = LoadBalancer
	installParams.UseCase = AllUseCases
	installParams.IngressInstallOption = StopIfInstalled
	*installParams.Domain = ""

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
        - name: DOMAIN
          value: 
        - name: INGRESS
          value: istio
        - name: USE_CASE
          value: all
        - name: INGRESS_INSTALL_OPTION
          value: StopIfInstalled
      restartPolicy: Never
`
	if res != expected {
		t.Error("installation manifest does not match")
	}
}

func resetFlagValues() {
	*installParams.ConfigFilePath = ""
	*installParams.InstallerImage = ""
	*installParams.PlatformIdentifier = "gke"
	*installParams.GatewayInput = "LoadBalancer"
	*installParams.Domain = ""
	*installParams.UseCaseInput = "all"
	*installParams.IngressInstallOptionInput = "StopIfInstalled"
}

func TestInstallCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("install --platform=gke --mock")

	resetFlagValues()

	r := newRedirector()
	r.redirectStdOut()

	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}

	out := r.revertStdOut()
	if !strings.Contains(out, "Used Installer version: docker.io/keptn/installer:latest") {
		t.Errorf("unexpected used version: %s", out)
	}
}

func TestInstallCmdWithKeptnVersion(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("install --platform=gke --keptn-installer-image=docker.io/keptn/installer:0.6.0 --mock")

	resetFlagValues()

	r := newRedirector()
	r.redirectStdOut()

	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}

	out := r.revertStdOut()
	if !strings.Contains(out, "Used Installer version: docker.io/keptn/installer:0.6.0") {
		t.Errorf("unexpected used version: %s", out)
	}
}

func TestInstallCmdWithGateway(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("install --platform=gke --gateway=NodePort --mock")

	resetFlagValues()

	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
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
        image: docker.io/keptn/installer:latest
        env:
        - name: PLATFORM
          value: gke
        - name: GATEWAY_TYPE
          value: NodePort
        - name: DOMAIN
          value: 
        - name: INGRESS
          value: istio
        - name: USE_CASE
          value: all
        - name: INGRESS_INSTALL_OPTION
          value: StopIfInstalled
      restartPolicy: Never
`
	if res != expected {
		t.Error("installation manifest does not match")
	}
}

func TestInstallCmdWithDomain(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("install --platform=gke --gateway=NodePort --domain=127.0.0.1.nip.io --mock")

	resetFlagValues()

	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
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
        image: docker.io/keptn/installer:latest
        env:
        - name: PLATFORM
          value: gke
        - name: GATEWAY_TYPE
          value: NodePort
        - name: DOMAIN
          value: 127.0.0.1.nip.io
        - name: INGRESS
          value: istio
        - name: USE_CASE
          value: all
        - name: INGRESS_INSTALL_OPTION
          value: StopIfInstalled
      restartPolicy: Never
`
	if res != expected {
		t.Error("installation manifest does not match")
	}
}

func TestInstallCmdWithUseCase(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("install --use-case=quality-gates --mock")

	resetFlagValues()

	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
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
        image: docker.io/keptn/installer:latest
        env:
        - name: PLATFORM
          value: gke
        - name: GATEWAY_TYPE
          value: LoadBalancer
        - name: DOMAIN
          value: 
        - name: INGRESS
          value: nginx
        - name: USE_CASE
          value: quality-gates
        - name: INGRESS_INSTALL_OPTION
          value: StopIfInstalled
      restartPolicy: Never
`
	if res != expected {
		t.Error("installation manifest does not match")
	}
}

func TestInstallCmdWithIstioInstallOption(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("install --ingress-install-option=Reuse --mock")

	resetFlagValues()

	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
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
        image: docker.io/keptn/installer:latest
        env:
        - name: PLATFORM
          value: gke
        - name: GATEWAY_TYPE
          value: LoadBalancer
        - name: DOMAIN
          value: 
        - name: INGRESS
          value: istio
        - name: USE_CASE
          value: all
        - name: INGRESS_INSTALL_OPTION
          value: Reuse
      restartPolicy: Never
`
	if res != expected {
		t.Error("installation manifest does not match")
	}
}
