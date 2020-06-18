package cmd

const installerJob = `---
apiVersion: batch/v1
kind: Job
metadata:
  name: installer
  namespace: keptn
spec:
  backoffLimit: 0
  template:
    metadata:
      labels:
        app: installer
    spec:
      volumes:
      - name: kubectl
        emptyDir: {}
      containers:
      - name: keptn-installer
        image: INSTALLER_IMAGE_PLACEHOLDER
        env:
        - name: PLATFORM
          value: PLATFORM_PLACEHOLDER
        - name: GATEWAY_TYPE
          value: GATEWAY_TYPE_PLACEHOLDER
        - name: DOMAIN
          value: DOMAIN_PLACEHOLDER
        - name: INGRESS
          value: INGRESS_PLACEHOLDER
        - name: USE_CASE
          value: USE_CASE_PLACEHOLDER
        - name: INGRESS_INSTALL_OPTION
          value: INGRESS_INSTALL_OPTION_PLACEHOLDER
      restartPolicy: Never
      serviceAccountName: keptn-installer
`

const installerServiceAccount = `---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: keptn-installer
  namespace: keptn
`

const installerClusterAdminClusterRoleBinding = `---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: installer-cluster-admin
  labels:
    "app": "keptn"
subjects:
  - kind: ServiceAccount
    name: keptn-installer
    namespace: keptn
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
`

func getGetRoleForNamespace(namespace string) string {
	return `---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: installer-get-in-namespace
  namespace: ` + namespace + `
  labels:
    "app": "keptn"
rules:
  - apiGroups:
      - "*"
    resources:
      - "*"
    verbs:
      - "get"
      - "list"
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: installer-get-in-namespace
  namespace: ` + namespace + `
  labels:
    "app": "keptn"
subjects:
  - kind: ServiceAccount
    name: keptn-installer
    namespace: keptn
roleRef:
  kind: Role
  name: installer-get-in-namespace
  apiGroup: rbac.authorization.k8s.io
`
}

func getAdminRoleForNamespace(namespace string) string {
	return `---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: installer-admin-namespace
  namespace: ` + namespace + `
  labels:
    "app": "keptn"
rules:
  - apiGroups:
      - "*"
    resources:
      - "*"
    verbs:
      - "*"
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: installer-admin-namespace
  namespace: ` + namespace + `
  labels:
    "app": "keptn"
subjects:
  - kind: ServiceAccount
    name: keptn-installer
    namespace: keptn
roleRef:
  kind: Role
  name: installer-admin-namespace
  apiGroup: rbac.authorization.k8s.io
`
}

const natsOperatorRoles = `---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: installer-nats-operator
  labels:
    "app": "keptn"
rules:
# Allow creating CRDs
- apiGroups:
  - apiextensions.k8s.io
  resources:
  - customresourcedefinitions
  verbs: ["get", "list", "create", "update", "watch"]

# Allow all actions on NATS Operator manager CRDs
- apiGroups:
  - nats.io
  resources:
  - natsclusters
  - natsserviceroles
  verbs: ["*"]

# Allowed actions on Pods
- apiGroups: [""]
  resources:
  - pods
  verbs: ["create", "watch", "get", "patch", "update", "delete", "list"]

# Allowed actions on Services
- apiGroups: [""]
  resources:
  - services
  verbs: ["create", "watch", "get", "patch", "update", "delete", "list"]

# Allowed actions on Secrets
- apiGroups: [""]
  resources:
  - secrets
  verbs: ["create", "watch", "get", "update", "delete", "list"]

# Allow all actions on some special subresources
- apiGroups: [""]
  resources:
  - pods/exec
  - pods/log
  - serviceaccounts/token
  - events
  verbs: ["*"]

# Allow listing Namespaces and ServiceAccounts
- apiGroups: [""]
  resources:
  - namespaces
  - serviceaccounts
  verbs: ["list", "get", "watch"]

# Allow actions on Endpoints
- apiGroups: [""]
  resources:
  - endpoints
  verbs: ["create", "watch", "get", "update", "delete", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: installer-nats-operator
  labels:
    "app": "keptn"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: installer-nats-operator
subjects:
- kind: ServiceAccount
  name: keptn-installer
  namespace: keptn
`

const natsOperatorServer = `---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: installer-nats-server
  labels:
    "app": "keptn"
rules:
- apiGroups: [""]
  resources:
  - nodes
  verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: installer-nats-server
  labels:
    "app": "keptn"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: installer-nats-server
subjects:
- kind: ServiceAccount
  name: keptn-installer
  namespace: keptn
`

const natsClusterRole = `---
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: installer-nats-create-clusterroles
  labels:
    "app": "keptn"
rules:  
  - apiGroups:
      - "rbac.authorization.k8s.io"
    resources:
      - clusterroles
      - clusterrolebindings
    verbs:
      - create
      - get
    resourceNames:
      - "nats-operator"
      - "nats-operator-binding"
      - "nats-server"
      - "nats-server-binding"
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: installer-nats-create-clusterroles
  labels:
    "app": "keptn"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: installer-nats-create-clusterroles
subjects:
  - kind: ServiceAccount
    name: keptn-installer
    namespace: keptn
`
