module github.com/keptn/keptn/remediation-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.1-compat.0.20200406125548-5337a2e806c4
	github.com/keptn/kubernetes-utils v0.0.0-20200401103501-ae44a5ee0656
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.4.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/helm v2.14.3+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190313205120-8b27c41bdbb1
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20190314002537-50662da99b70
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190311093542-50b561225d70
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20190314002815-ce92c5cfdd61
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20190314002447-97ed623e3835
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20190314002154-4d735c31b054
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20190314002350-b74e9e79538d
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20190314002251-f6da02f58325
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20190314000836-236f85ce49e5
)
