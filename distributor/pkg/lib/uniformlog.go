package lib

import (
	"bufio"
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/api/utils"
	logger "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
	"time"
)

var errorPhrases = []string{"error", "could not", "fail", "couldn't"}

type UniformLog interface {
	Log(keptnapimodels.LogEntry) error
	Start(ctx context.Context)
}

type K8sUniformLogger struct {
	K8sClient    kubernetes.Interface
	Integration  keptnapimodels.Integration
	Closed       chan struct{}
	logCache     []keptnapimodels.LogEntry
	logHandler   keptn.LogHandler
	theClock     clock.Clock
	syncInterval time.Duration
}

func NewK8sUniformLogger(k8sClient kubernetes.Interface, integration keptnapimodels.Integration, theClock clock.Clock, syncInterval time.Duration) *K8sUniformLogger {
	return &K8sUniformLogger{
		K8sClient:    k8sClient,
		Integration:  integration,
		Closed:       make(chan struct{}),
		logCache:     []keptnapimodels.LogEntry{},
		logHandler:   keptn.LogHandler{},
		theClock:     theClock,
		syncInterval: syncInterval,
	}
}

func (K8sUniformLogger) Log(entry keptnapimodels.LogEntry) error {
	panic("implement me")
}

func (kl *K8sUniformLogger) Start(ctx context.Context) {
	go func() {
		rs, err := kl.K8sClient.CoreV1().Pods(kl.Integration.MetaData.KubernetesMetaData.Namespace).GetLogs(kl.Integration.MetaData.KubernetesMetaData.PodName, &v1.PodLogOptions{
			Container:  kl.Integration.MetaData.KubernetesMetaData.DeploymentName,
			Follow:     true,
			Timestamps: false,
		}).Stream(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer rs.Close()

		go func() {
			<-kl.Closed
			rs.Close()
		}()

		sc := bufio.NewScanner(rs)

		for sc.Scan() {
			if containsErrorMessage(sc.Text()) {
				kl.logCache = append(kl.logCache, keptnapimodels.LogEntry{
					IntegrationID: kl.Integration.ID,
					Message:       sc.Text(),
				})
			}
		}
	}()

	go func() {
		ticker := kl.theClock.Ticker(kl.syncInterval)
		for {
			select {
			case <-ctx.Done():
				logger.Info("cancelling event dispatcher loop")
				return
			case <-ticker.C:
				logger.Debugf("%.2f seconds have passed. Sending log messages to API", kl.syncInterval.Seconds())
				kl.logHandler.Log(kl.logCache)
			}
		}
	}()

	go func() {
		<-ctx.Done()
		close(kl.Closed)
	}()
}

func containsErrorMessage(text string) bool {
	myText := strings.ToLower(text)
	for _, phrase := range errorPhrases {
		if strings.Contains(myText, phrase) {
			return true
		}
	}
	return false
}
