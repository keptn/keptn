package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	logger "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

var deploymentsToCheck = []string{
	// core
	"approval-service",
	"lighthouse-service",
	"remediation-service",
	"secret-service",
	"statistics-service",
	"webhook-service",
	// execution plane
	"helm-service",
	"argo-service",
	"jmeter-service",
	"job-executor-service",
}

var lastEventReceived time.Time

func EnsureDeploymentsAreUp() (bool, error) {
	// reset lastEventReceived
	lastEventReceived = time.Now()

	// get k8s clients
	var k8sClient kubernetes.Interface

	var namespace string
	namespace = os.Getenv("POD_NAMESPACE")

	var deploymentsToScaleUp []string

	k8sClient, err := keptnutils.GetClientset(true)
	if err != nil {
		logger.Debug("Could not create kubernetes client, will skip checking k8s deployments.")
		// we need to assume that all deployments are up - but we will also return the error
		return true, err
	}

	deploymentsClient := k8sClient.AppsV1().Deployments(namespace)

	for _, deploymentName := range deploymentsToCheck {
		// check if deployment exists
		deployment, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})

		if err != nil {
			// probably failed to get deployment, means we can ignore it and continue with the next one
			logger.Debugf("Couldn't find deployment %s, skipping it...", deploymentName)
			continue
		}

		// else: ensure the deployment is scaled up
		if deployment.Status.AvailableReplicas == 0 {
			logger.Infof("Deployment %s is not up... Scaling it up now!", deploymentName)

			deploymentsToScaleUp = append(deploymentsToScaleUp, deploymentName)

			// try to scale up
			err = keptnutils.ScaleDeployment(true, deploymentName, namespace, 1)

			if err != nil {
				logger.Errorf("Failed to update deployment: %s", err.Error())
				return false, errors.New(fmt.Sprintf("Deployment %s is not up", deploymentName))
			}
		} else {
			logger.Debugf("Deployment %s is up and running, continuing...", deploymentName)
		}
	}

	// check if all deployments have been scaled up
	for _, deploymentName := range deploymentsToScaleUp {
		// wait for deployment is up
		logger.Debugf("Checking if deployment %s is up...", deploymentName)
		if err := keptnutils.WaitForDeploymentToBeRolledOut(true, deploymentName, namespace); err != nil {
			return false, fmt.Errorf("Error when waiting for deployment %s in namespace %s: %s", deploymentName, namespace, err.Error())
		}
	}

	if len(deploymentsToScaleUp) > 0 {
		// add an extra sleep to be 100% sure
		time.Sleep(5 * time.Second)
	}

	logger.Infof("All deployments are up!")

	return true, nil
}

func SailDownLoop(ctx context.Context) {
	// init last event received time
	lastEventReceived = time.Now()

	// check if we should scale down
	for {
		if lastEventReceived.Add(2 * time.Minute).Before(time.Now()) {
			logger.Infof("too much time since last cloud-event... saildown now!")
			scaleDownDeployments()

			// no need to check for another 2 minutes
			time.Sleep(2 * time.Minute)
		}
		// check again in 30 seconds
		time.Sleep(30 * time.Second)
	}
}

func scaleDownDeployments() {
	// get k8s clients
	var k8sClient kubernetes.Interface

	var namespace string
	namespace = os.Getenv("POD_NAMESPACE")

	k8sClient, err := keptnutils.GetClientset(true)
	if err != nil {
		logger.Debug("Could not create kubernetes client, will skip checking k8s deployments.")
		// we need to assume that all deployments are up - but we will also return the error
		return
	}

	deploymentsClient := k8sClient.AppsV1().Deployments(namespace)

	for _, deploymentName := range deploymentsToCheck {
		// check if deployment exists
		deployment, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})

		if err != nil {
			// probably failed to get deployment, means we can ignore it and continue with the next one
			logger.Debugf("Couldn't find deployment %s, skipping it...", deploymentName)
			continue
		}

		// else: ensure the deployment is scaled up
		if deployment.Status.AvailableReplicas > 0 {
			// scale it down
			err = keptnutils.ScaleDeployment(true, deploymentName, namespace, 0)

			if err != nil {
				logger.Errorf("Failed to scale down deployment: %s", err.Error())
			}
		}
	}
}
