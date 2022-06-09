package go_tests

import (
	"context"
	"fmt"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"k8s.io/apimachinery/pkg/util/wait"
	"net"
	"os/exec"
	"time"
)

const kubectlExecutable = "kubectl"

// KubeCtlApplyFromURL wraps the kubectl command line tool in order to perform a "kubectl apply" command
// with resources downloaded from the given "url". The default namespace (set via KEPTN_NAMESPACE) will be used
// but can be overridden using the "namespace" param. The function returns a function which can be called to
// apply the corresponding "kubectl delete" command to undo the "kubectl apply" command
func KubeCtlApplyFromURL(url string, namespace ...string) (func() error, error) {
	var ns = GetKeptnNameSpaceFromEnv()
	if len(namespace) == 1 {
		ns = namespace[0]
	}
	fmt.Printf("Executing: %s %s -n=%s -f=%s\n", kubectlExecutable, "apply", ns, url)
	result, err := keptnkubeutils.ExecuteCommand(kubectlExecutable, []string{"apply", "-n=" + ns, "-f=" + url})
	if err != nil {
		return nil, err
	}
	fmt.Println(result)

	deleteFunc := func() error {
		return KubeCtlDeleteFromURL(url, namespace...)
	}
	return deleteFunc, err
}

func KubeCtlDeleteFromURL(url string, namespace ...string) error {
	var ns = GetKeptnNameSpaceFromEnv()
	if len(namespace) == 1 {
		ns = namespace[0]
	}
	fmt.Printf("Executing: %s %s -n=%s -f=%s\n", kubectlExecutable, "delete", ns, url)
	result, err := keptnkubeutils.ExecuteCommand(kubectlExecutable, []string{"delete", "-n=" + ns, "-f=" + url})
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func KubeCtlPortForwardSvc(ctx context.Context, svcName, localPort string, remotePort string, namespace ...string) error {
	var ns = GetKeptnNameSpaceFromEnv()
	if len(namespace) == 1 {
		ns = namespace[0]
	}
	fmt.Printf("Executing: %s port-forward -n %s %s %s\n", kubectlExecutable, ns, svcName, localPort+":"+remotePort)
	cmd := exec.CommandContext(ctx, kubectlExecutable, "port-forward", "-n", ns, svcName, localPort+":"+remotePort)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = wait.Poll(2*time.Second, 30*time.Second, func() (bool, error) {
		_, err := net.DialTimeout("tcp", "127.0.0.1:"+localPort, 2*time.Second)
		return err == nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}
