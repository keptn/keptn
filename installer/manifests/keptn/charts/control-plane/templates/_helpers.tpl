{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "control-plane.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "control-plane.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "control-plane.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "control-plane.labels" -}}
helm.sh/chart: {{ include "control-plane.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Values.global.keptn.appVersion | default .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "control-plane.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "control-plane.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "control-plane.dist.livenessProbe" -}}
livenessProbe:
  httpGet:
    path: /health
    port: {{.port | default 8080}}
  initialDelaySeconds: {{.initialDelaySeconds | default 10}}
  periodSeconds: 5
{{- end }}

{{- define "control-plane.dist.readinessProbe" -}}
readinessProbe:
  httpGet:
    path: /health
    port: {{.port | default 8080}}
  initialDelaySeconds: {{.initialDelaySeconds | default 5}}
  periodSeconds: 5
{{- end }}

{{/*
preStop hook for control-plane deployments
*/}}
{{- define "control-plane.prestop" -}}
lifecycle:
  preStop:
    exec:
      # using 90s of sleeping to be on the safe side before terminating the pod
      command: ["/bin/sleep", {{ . }} ]
{{- end }}

{{- define "control-plane.dist.common.env.vars" -}}
- name: PUBSUB_URL
  value: 'nats://keptn-nats'
- name: VERSION
  valueFrom:
    fieldRef:
      fieldPath: metadata.labels['app.kubernetes.io/version']
- name: DISTRIBUTOR_VERSION
  value: {{ include "service.tag" ( dict "imageRoot" .Values.distributor.image "global" .Values.global.keptn "defaultTag" .Chart.AppVersion) | quote}}
- name: LOCATION
  valueFrom:
   fieldRef:
      fieldPath: metadata.labels['app.kubernetes.io/component']
- name: K8S_DEPLOYMENT_NAME
  valueFrom:
    fieldRef:
      fieldPath: metadata.labels['app.kubernetes.io/name']
- name: K8S_POD_NAME
  valueFrom:
    fieldRef:
     fieldPath: metadata.name
- name: K8S_NAMESPACE
{{- if .Values.distributor.metadata.namespace }}
  value: {{ .Values.distributor.metadata.namespace }}
{{- else }}
  valueFrom:
    fieldRef:
      fieldPath: metadata.namespace
{{- end }}
- name: K8S_NODE_NAME
{{- if .Values.distributor.metadata.hostname }}
  value: {{ .Values.distributor.metadata.hostname }}
{{- else }}
  valueFrom:
    fieldRef:
      fieldPath: spec.nodeName
{{- end }}
{{- if .Values.distributor.config.queueGroup.enabled }}
- name: PUBSUB_GROUP
  valueFrom:
    fieldRef:
      fieldPath: metadata.labels['app.kubernetes.io/name']
{{- end }}
- name: OAUTH_CLIENT_ID
  value: "{{ (((.Values.distributor).config).oauth).clientID }}"
- name: OAUTH_CLIENT_SECRET
  value: "{{ (((.Values.distributor).config).oauth).clientSecret }}"
- name: OAUTH_DISCOVERY
  value: "{{ (((.Values.distributor).config).oauth).discovery }}"
- name: OAUTH_TOKEN_URL
  value: "{{ (((.Values.distributor).config).oauth).tokenURL }}"
- name: OAUTH_SCOPES
  value: "{{ (((.Values.distributor).config).oauth).scopes }}"
{{- end }}

{{- define "control-plane.common.security-context-seccomp" -}}
{{- if ge .Capabilities.KubeVersion.Minor "21" }}
  seccompProfile:
    type: RuntimeDefault
{{- end -}}
{{- end -}}

{{- define "control-plane.bridge.pod-security-context" -}}
{{- if .Values.bridge.podSecurityContext -}}
{{- if .Values.bridge.podSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.bridge.podSecurityContext "enabled" "defaultSeccompProfile" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- if not .Values.bridge.podSecurityContext.seccompProfile }}
{{- if .Values.bridge.podSecurityContext.defaultSeccompProfile -}}
{{- include "control-plane.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  fsGroup: 65532
{{- include "control-plane.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}

{{- define "control-plane.bridge.container-security-context" -}}
{{- if .Values.bridge.containerSecurityContext -}}
{{- if .Values.bridge.containerSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.bridge.containerSecurityContext "enabled" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  runAsNonRoot: true
  runAsUser: 65532
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
  privileged: false
{{- end -}}
{{- end -}}

{{- define "control-plane.apiGatewayNginx.pod-security-context" -}}
{{- if .Values.apiGatewayNginx.podSecurityContext -}}
{{- if .Values.apiGatewayNginx.podSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.apiGatewayNginx.podSecurityContext "enabled" "defaultSeccompProfile" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- if not .Values.apiGatewayNginx.podSecurityContext.seccompProfile -}}
{{- if .Values.apiGatewayNginx.podSecurityContext.defaultSeccompProfile -}}
{{- include "control-plane.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  fsGroup: 101
{{- include "control-plane.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}

{{- define "control-plane.apiGatewayNginx.container-security-context" -}}
{{- if .Values.apiGatewayNginx.containerSecurityContext -}}
{{- if .Values.apiGatewayNginx.containerSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.apiGatewayNginx.containerSecurityContext "enabled" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  runAsNonRoot: true
  runAsUser: 101
  readOnlyRootFilesystem: false
  allowPrivilegeEscalation: false
  privileged: false
{{- end -}}
{{- end -}}

{{- define "control-plane.common.pod-security-context" -}}
{{- if (.Values.common).podSecurityContext -}}
{{- if .Values.common.podSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.common.podSecurityContext "enabled" "defaultSeccompProfile" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- if not .Values.common.podSecurityContext.seccompProfile -}}
{{- if .Values.apiGatewayNginx.podSecurityContext.defaultSeccompProfile -}}
{{- include "control-plane.common.security-context-seccomp" . -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  fsGroup: 65532
{{- include "control-plane.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}

{{- define "control-plane.common.container-security-context" -}}
{{- if (.Values.common).containerSecurityContext -}}
{{- if .Values.common.containerSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.common.containerSecurityContext "enabled" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  runAsNonRoot: true
  runAsUser: 65532
  readOnlyRootFilesystem: false
  allowPrivilegeEscalation: false
  privileged: false
{{- end -}}
{{- end -}}

{{/*
rollingUpdate upgrade strategy for control plane deployments
*/}}
{{- define "control-plane.common.update-strategy" -}}
{{- if (.Values.common).strategy -}}
strategy:
{{- toYaml .Values.common.strategy | nindent 2 -}}
{{- else -}}
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1
    maxUnavailable: 0
{{- end -}}
{{- end -}}


