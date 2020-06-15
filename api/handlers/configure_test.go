package handlers

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"

	networking "k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAddBridgeToIngress(t *testing.T) {
	ingress := &networking.Ingress{
		TypeMeta: v1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		Spec: networking.IngressSpec{
			TLS: []networking.IngressTLS{
				{
					Hosts:      []string{"'*.keptn.test.domain'"},
					SecretName: "sslcerts",
				},
			},
			Rules: []networking.IngressRule{
				{
					Host: "api.keptn.test.domain",
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Backend: networking.IngressBackend{
										ServiceName: "api-gateway-nginx",
										ServicePort: intstr.IntOrString{IntVal: 80},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	addBridgeToIngress("test.domain", ingress)
	if len(ingress.Spec.Rules) != 2 {
		t.Error("Unexpected number of rules")
	}
	if ingress.Spec.Rules[1].Host != "bridge.keptn.test.domain" {
		t.Error("Unexpected name of host")
	}
}

func TestAddBridgeToIngressWithExistingBridgeHost(t *testing.T) {
	ingress := &networking.Ingress{
		TypeMeta: v1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		Spec: networking.IngressSpec{
			TLS: []networking.IngressTLS{
				{
					Hosts:      []string{"'*.keptn.test.domain'"},
					SecretName: "sslcerts",
				},
			},
			Rules: []networking.IngressRule{
				{
					Host: "api.keptn.test.domain",
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Backend: networking.IngressBackend{
										ServiceName: "api-gateway-nginx",
										ServicePort: intstr.IntOrString{IntVal: 80},
									},
								},
							},
						},
					},
				},
				getBridgeRule("test.domain"),
			},
		},
	}

	addBridgeToIngress("test.domain", ingress)
	if len(ingress.Spec.Rules) != 2 {
		t.Error("Unexpected number of rules")
	}
	if ingress.Spec.Rules[1].Host != "bridge.keptn.test.domain" {
		t.Error("Unexpected name of host")
	}
}

func TestDisposeBridgeFromIngressWithNoHost(t *testing.T) {
	ingress := &networking.Ingress{
		TypeMeta: v1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		Spec: networking.IngressSpec{
			TLS: []networking.IngressTLS{
				{
					Hosts:      []string{"'*.keptn.test.domain'"},
					SecretName: "sslcerts",
				},
			},
			Rules: []networking.IngressRule{
				{
					Host: "api.keptn.test.domain",
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Backend: networking.IngressBackend{
										ServiceName: "api-gateway-nginx",
										ServicePort: intstr.IntOrString{IntVal: 80},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	removeBridgeFromIngress(ingress)
	if len(ingress.Spec.Rules) != 1 {
		t.Error("Unexpected number of rules")
	}
	if ingress.Spec.Rules[0].Host != "api.keptn.test.domain" {
		t.Error("Unexpected name of host")
	}
}

func Test_getHostForBridge(t *testing.T) {
	type args struct {
		keptnDomain string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get bridge hostname",
			args: args{
				keptnDomain: "my-domain.com",
			},
			want: "bridge.keptn.my-domain.com",
		},
		{
			name: "get bridge hostname from domain containing a port",
			args: args{
				keptnDomain: "my-domain.com:1234",
			},
			want: "bridge.keptn.my-domain.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHostForBridge(tt.args.keptnDomain); got != tt.want {
				t.Errorf("getHostForBridge() = %v, want %v", got, tt.want)
			}
		})
	}
}
