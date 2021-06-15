package go_tests

import (
	"fmt"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
)

const kubectExecutable = "kubectl"

// KubeCtlApplyFromURL wraps the kubectl command line tool in order to perform a "kubectl apply" command
// with resources downloaded from the given "url". The default namespace (set via KEPTN_NAMESPACE) will be used
// but can be overridden using the "namespace" param. The function returns a function which can be called to
// apply the corresponding "kubectl delete" command to undo the "kubectl apply" command
func KubeCtlApplyFromURL(url string, namespace ...string) (func() error, error) {
	var ns = GetKeptnNameSpaceFromEnv()
	if len(namespace) == 1 {
		ns = namespace[0]
	}
	fmt.Printf("Executing: %s %s -n=%s -f=%s\n", kubectExecutable, "apply", ns, url)
	result, err := keptnkubeutils.ExecuteCommand(kubectExecutable, []string{"apply", "-n=" + ns, "-f=" + url})
	if err != nil {
		return nil, err
	}
	fmt.Println(result)

	deleteFunc := func() error {
		var ns = GetKeptnNameSpaceFromEnv()
		if len(namespace) == 1 {
			ns = namespace[0]
		}
		fmt.Printf("Executing: %s %s -n=%s -f=%s\n", kubectExecutable, "delete", ns, url)
		result, err = keptnkubeutils.ExecuteCommand(kubectExecutable, []string{"delete", "-n=" + ns, "-f=" + url})
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	}
	return deleteFunc, err
}
