module github.com/keptn/keptn/remediation-service

go 1.12

require (
	github.com/DATA-DOG/go-sqlmock v1.3.3 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/chai2010/gettext-go v0.0.0-20170215093142-bf70f2a70fb1 // indirect
	github.com/cloudevents/sdk-go v0.10.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/keptn/go-utils v0.3.1-0.20191111100301-bff0cac85494
	github.com/lib/pq v1.2.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/rubenv/sql-migrate v0.0.0-20191022111038-5cdff0d8cc42 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20191004055002-72853e10c5a3 // indirect
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/square/go-jose.v2 v2.4.0 // indirect
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0
	k8s.io/apiextensions-apiserver v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/apimachinery v0.0.0
	k8s.io/apiserver v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/cli-runtime v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cloud-provider v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/helm v2.14.3+incompatible
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d // indirect
	k8s.io/kubernetes v1.14.0
	sigs.k8s.io/kustomize v2.0.3+incompatible // indirect
	sigs.k8s.io/yaml v1.1.0
	vbom.ml/util v0.0.0-20180919145318-efcd4e0f9787 // indirect
)

replace (
	github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
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
	//k8s.io/cri-api => k8s.io/cri-api kubernetes-1.14.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20190314002815-ce92c5cfdd61
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20190314002447-97ed623e3835
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20190314002154-4d735c31b054
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20190314002350-b74e9e79538d
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20190314002251-f6da02f58325
	//k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers kubernetes-1.14.0
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20190314000836-236f85ce49e5
)
