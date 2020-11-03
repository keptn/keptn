package common

import "fmt"

func printQuickAccessInstructionsForUnix(keptnNamespace, keptnReleaseDocsURL string, endpoint string) {
	if endpoint == "http://localhost:8080/api" {
		fmt.Println("* To quickly access Keptn, you can use a port-forward and then authenticate your Keptn CLI:\n" +
			" - kubectl -n " + keptnNamespace + " port-forward service/api-gateway-nginx 8080:80\n" +
			" - keptn auth --endpoint=http://localhost:8080/api --api-token=$(kubectl get secret keptn-api-token -n " + keptnNamespace + " -ojsonpath={.data.keptn-api-token} | base64 --decode)\n")
	} else {
		fmt.Println("* To quickly access Keptn, you can authenticate your Keptn CLI (in a Linux shell):\n" +
			" - keptn auth --endpoint=" + endpoint + " --api-token=$(kubectl get secret keptn-api-token -n " + keptnNamespace + " -ojsonpath={.data.keptn-api-token} | base64 --decode)\n")
	}

	fmt.Println("* Alternatively, follow the instructions provided at: https://keptn.sh/docs/" + keptnReleaseDocsURL + "/operate/install/#authenticate-keptn-cli\n")

	printExposeEndpointInstructions(keptnReleaseDocsURL)
}

func printQuickAccessInstructionsForWindows(keptnReleaseDocsURL string) {
	fmt.Println("* To quickly access Keptn, you can use a port-forward and then authenticate your Keptn CLI as described here: https://keptn.sh/docs/" + keptnReleaseDocsURL + "/operate/install/#authenticate-keptn-cli\n")

	printExposeEndpointInstructions(keptnReleaseDocsURL)
}

func printExposeEndpointInstructions(keptnReleaseDocsURL string) {
	fmt.Println("* To expose Keptn on a public endpoint, please continue with the installation guidelines provided at:\n" +
		" - https://keptn.sh/docs/" + keptnReleaseDocsURL + "/operate/install#install-keptn\n")
}
