package common

func PrintQuickAccessInstructions(keptnNamespace, keptnReleaseDocsURL string) {
	fmt.Println("* * To quickly access Keptn, you can use a port-forward and then authenticate your Keptn CLI as described here: https://keptn.sh/docs/" + keptnReleaseDocsURL + "/operate/install/#authenticate-keptn-cli\n")
	fmt.Println("* To expose Keptn on a public endpoint, please continue with the installation guidelines provided at:\n" +
		" - https://keptn.sh/docs/" + keptnReleaseDocsURL + "/operate/install#install-keptn\n")
}
