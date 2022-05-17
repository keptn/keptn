package zero_downtime

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"

	testutils "github.com/keptn/keptn/test/go-tests"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"testing"
)

func TestEvaluationsWithApproval(t *testing.T) {
	images := []string{"0.15.1-dev.202205161008", "0.15.1-dev.202205161146"}
	lhImages := []string{"ready1", "ready2"}
	services := []string{"api", "resource-service", "lighthouse-service"}
	//services := []string{"lighthouse-service"}

	ctx, cancel := context.WithCancel(context.Background())

	for _, svc := range services {
		go func(service string) {
			imgArr := images
			if service == "lighthouse-service" {
				imgArr = lhImages
			}
			err := updateImageOfService(ctx, t, service, imgArr)
			if err != nil {
				t.Logf("%v", err)
			}
		}(svc)
	}

	doEvaluations()
	//<-time.After(2 * time.Minute)
	cancel()
}

func doEvaluations() {
	for i := 0; i < 100; i++ {
		nrEvaluations := 0
		go func() {
			_, err := triggerEvaluation("podtatohead", "hardening", "helloservice")
			if err != nil {
				nrEvaluations++
			}
		}()

		<-time.After(5 * time.Second)
	}
}

func triggerEvaluation(projectName, stageName, serviceName string) (string, error) {
	cliResp, err := testutils.ExecuteCommand(fmt.Sprintf("keptn trigger evaluation --project=%s --stage=%s --service=%s --timeframe=5m", projectName, stageName, serviceName))

	if err != nil {
		return "", err
	}
	var keptnContext string
	split := strings.Split(cliResp, "\n")
	for _, line := range split {
		if strings.Contains(line, "ID of") {
			splitLine := strings.Split(line, ":")
			if len(splitLine) == 2 {
				keptnContext = strings.TrimSpace(splitLine[1])
			}
		}
	}
	return keptnContext, err
}

func updateImageOfService(ctx context.Context, t *testing.T, service string, images []string) error {
	clientset, err := keptnkubeutils.GetClientset(false)

	if err != nil {
		return err
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			nextImage := images[i%len(images)]
			get, err := clientset.AppsV1().Deployments(testutils.GetKeptnNameSpaceFromEnv()).Get(context.TODO(), service, v1.GetOptions{})
			if err != nil {
				break
			}

			imageWithTag := get.Spec.Template.Spec.Containers[0].Image
			split := strings.Split(imageWithTag, ":")
			updatedImage := fmt.Sprintf("%s:%s", split[0], nextImage)

			get.Spec.Template.Spec.Containers[0].Image = updatedImage

			t.Logf("upgrading %s to %s", service, updatedImage)
			_, err = clientset.AppsV1().Deployments(testutils.GetKeptnNameSpaceFromEnv()).Update(context.TODO(), get, v1.UpdateOptions{})
			if err != nil {
				break
			}

			require.Eventually(t, func() bool {
				pods, err := clientset.CoreV1().Pods(testutils.GetKeptnNameSpaceFromEnv()).List(context.TODO(), v1.ListOptions{LabelSelector: "app.kubernetes.io/name=" + service})
				if err != nil {
					return false
				}

				if int32(len(pods.Items)) != 1 {
					// make sure only one pod is running
					return false
				}

				for _, item := range pods.Items {
					if len(item.Spec.Containers) == 0 {
						continue
					}
					if item.Spec.Containers[0].Image == updatedImage {
						return true
					}
				}
				return false
			}, 3*time.Minute, 10*time.Second)
			<-time.After(5 * time.Second)
			i++
		}
	}
}
