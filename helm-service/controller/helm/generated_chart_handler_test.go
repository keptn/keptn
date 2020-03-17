package helm

import (
	"reflect"
	"testing"

	"github.com/keptn/keptn/helm-service/controller/mesh"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func TestGeneratedChartHandler_generateServices(t *testing.T) {
	type fields struct {
		mesh        mesh.Mesh
		keptnDomain string
	}
	type args struct {
		svc       *corev1.Service
		project   string
		stageName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*chart.Template
		wantErr bool
	}{
		{
			name: "Create a service",
			fields: fields{
				mesh:        mesh.NewIstioMesh(),
				keptnDomain: "keptn.com",
			},
			args: args{
				svc: &corev1.Service{
					ObjectMeta: v1.ObjectMeta{
						Name: "carts",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "carts",
						},
					},
				},
				project:   "sockshop",
				stageName: "staging",
			},
			want: []*chart.Template{
				{
					Name:                 "templates/carts-canary-service.yaml",
					Data:                 []byte("apiVersion: v1\nkind: Service\nmetadata:\n  creationTimestamp: null\n  name: carts-canary\nspec:\n  selector:\n    app: carts\nstatus:\n  loadBalancer: {}\n"),
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
				{
					Name:                 "templates/carts-canary-istio-destinationrule.yaml",
					Data:                 []byte("apiVersion: networking.istio.io/v1alpha3\nkind: DestinationRule\nmetadata:\n  creationTimestamp: null\n  name: carts-canary\nspec:\n  host: carts-canary.sockshop-staging.svc.cluster.local\n"),
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
				{
					Name:                 "templates/carts-primary-service.yaml",
					Data:                 []byte("apiVersion: v1\nkind: Service\nmetadata:\n  creationTimestamp: null\n  name: carts-primary\nspec:\n  selector:\n    app: carts-primary\nstatus:\n  loadBalancer: {}\n"),
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
				{
					Name:                 "templates/carts-primary-istio-destinationrule.yaml",
					Data:                 []byte("apiVersion: networking.istio.io/v1alpha3\nkind: DestinationRule\nmetadata:\n  creationTimestamp: null\n  name: carts-primary\nspec:\n  host: carts-primary.sockshop-staging.svc.cluster.local\n"),
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
				{
					Name:                 "templates/carts-istio-virtualservice.yaml",
					Data:                 []byte("apiVersion: networking.istio.io/v1alpha3\nkind: VirtualService\nmetadata:\n  creationTimestamp: null\n  name: carts\nspec:\n  gateways:\n  - public-gateway.istio-system\n  - mesh\n  hosts:\n  - carts.sockshop-staging.keptn.com\n  - carts\n  http:\n  - route:\n    - destination:\n        host: carts-canary.sockshop-staging.svc.cluster.local\n    - destination:\n        host: carts-primary.sockshop-staging.svc.cluster.local\n      weight: 100\n"),
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GeneratedChartHandler{
				mesh:        tt.fields.mesh,
				keptnDomain: tt.fields.keptnDomain,
			}
			got, err := c.generateServices(tt.args.svc, tt.args.project, tt.args.stageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateServices() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratedChartHandler_generateDeployment(t *testing.T) {
	type fields struct {
		mesh        mesh.Mesh
		keptnDomain string
	}
	type args struct {
		depl *appsv1.Deployment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *chart.Template
		wantErr bool
	}{
		{
			name: "Create deployment",
			fields: fields{
				mesh:        mesh.NewIstioMesh(),
				keptnDomain: "keptn.com",
			},
			args: args{
				depl: &appsv1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name: "carts",
					},
					Spec: appsv1.DeploymentSpec{
						Selector: &v1.LabelSelector{
							MatchLabels: map[string]string{
								"app.kubernetes.io/name": "carts",
							},
						},
						Template: corev1.PodTemplateSpec{},
					},
				},
			},
			want: &chart.Template{
				Name:                 "templates/carts-primary-deployment.yaml",
				Data:                 []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  creationTimestamp: null\n  name: carts-primary\nspec:\n  selector:\n    matchLabels:\n      app.kubernetes.io/name: carts-primary\n  strategy: {}\n  template:\n    metadata:\n      creationTimestamp: null\n    spec:\n      containers: null\nstatus: {}\n"),
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     nil,
				XXX_sizecache:        0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GeneratedChartHandler{
				mesh:        tt.fields.mesh,
				keptnDomain: tt.fields.keptnDomain,
			}
			got, err := c.generateDeployment(tt.args.depl)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateDeployment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
