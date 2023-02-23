{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "keptn.name" -}}
{{- include "keptn.common.names.name" . -}}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "keptn.fullname" -}}
{{- include "keptn.common.names.fullname" . -}}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "keptn.chart" -}}
{{- include "keptn.common.names.chart" . -}}
{{- end }}

{{- define "keptn.dist.livenessProbe" -}}
livenessProbe:
  httpGet:
    path: /health
    port: {{.port | default 8080}}
  initialDelaySeconds: {{ .initialDelaySeconds | default 10 }}
  periodSeconds: 5
{{- end }}

{{- define "keptn.dist.readinessProbe" -}}
readinessProbe:
  httpGet:
    path: /health
    port: {{.port | default 8080}}
  initialDelaySeconds: {{ .initialDelaySeconds | default 5 }}
  periodSeconds: 5
{{- end }}

{{/*
preStop hook for keptn deployments
*/}}
{{- define "keptn.prestop" -}}
lifecycle:
  preStop:
    exec:
      # using 90s of sleeping to be on the safe side before terminating the pod
      command: ["/bin/sleep", {{ . }} ]
{{- end }}

{{- define "keptn.common.env.vars" -}}
- name: K8S_DEPLOYMENT_NAME
  valueFrom:
    fieldRef:
      apiVersion: v1
      fieldPath: 'metadata.labels[''app.kubernetes.io/name'']'
- name: K8S_DEPLOYMENT_VERSION
  valueFrom:
    fieldRef:
      apiVersion: v1
      fieldPath: 'metadata.labels[''app.kubernetes.io/version'']'
- name: K8S_NAMESPACE
{{- if .Values.distributor.metadata.namespace }}
  value: {{ .Values.distributor.metadata.namespace }}
{{- else }}
  valueFrom:
    fieldRef:
      apiVersion: v1
      fieldPath: metadata.namespace
{{- end }}
- name: K8S_NODE_NAME
{{- if .Values.distributor.metadata.hostname }}
  value: {{ .Values.distributor.metadata.hostname }}
{{- else }}
  valueFrom:
    fieldRef:
      apiVersion: v1
      fieldPath: spec.nodeName
{{- end }}
- name: K8S_POD_NAME
  valueFrom:
    fieldRef:
      apiVersion: v1
      fieldPath: metadata.name
{{- end }}

{{- define "keptn.dist.common.env.vars" -}}
- name: PUBSUB_URL
  value: 'nats://keptn-nats'
- name: VERSION
  valueFrom:
    fieldRef:
      fieldPath: 'metadata.labels[''app.kubernetes.io/version'']'
- name: DISTRIBUTOR_VERSION
  value: {{ include "keptn.common.images.tag" ( dict "imageRoot" .Values.distributor.image "global" .Values.global.keptn "defaultTag" .Chart.AppVersion) | quote }}
- name: API_PROXY_HTTP_TIMEOUT
  value: {{ ((.Values.distributor.config).proxy).httpTimeout | default "30" | quote }}
- name: API_PROXY_MAX_PAYLOAD_BYTES_KB
  value: {{ ((.Values.distributor.config).proxy).maxPayloadBytesKB | default "64" | quote }}
- name: K8S_DEPLOYMENT_NAME
  valueFrom:
    fieldRef:
      fieldPath: 'metadata.labels[''app.kubernetes.io/name'']'
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
      fieldPath: 'metadata.labels[''app.kubernetes.io/name'']'
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

{{- define "keptn.common.security-context-seccomp" -}}
{{- if ge .Capabilities.KubeVersion.Minor "21" }}
  seccompProfile:
    type: RuntimeDefault
{{- end -}}
{{- end -}}

{{- define "keptn.bridge.pod-security-context" -}}
{{- if .Values.bridge.podSecurityContext -}}
{{- if .Values.bridge.podSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.bridge.podSecurityContext "enabled" "defaultSeccompProfile" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- if not .Values.bridge.podSecurityContext.seccompProfile }}
{{- if .Values.bridge.podSecurityContext.defaultSeccompProfile -}}
{{- include "keptn.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  fsGroup: 65532
{{- include "keptn.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}

{{- define "keptn.bridge.container-security-context" -}}
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

{{- define "keptn.apiGatewayNginx.pod-security-context" -}}
{{- if .Values.apiGatewayNginx.podSecurityContext -}}
{{- if .Values.apiGatewayNginx.podSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.apiGatewayNginx.podSecurityContext "enabled" "defaultSeccompProfile" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- if not .Values.apiGatewayNginx.podSecurityContext.seccompProfile -}}
{{- if .Values.apiGatewayNginx.podSecurityContext.defaultSeccompProfile -}}
{{- include "keptn.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  fsGroup: 101
{{- include "keptn.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}

{{- define "keptn.apiGatewayNginx.container-security-context" -}}
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

{{- define "keptn.common.pod-security-context" -}}
{{- if .Values.podSecurityContext -}}
{{- if .Values.podSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.podSecurityContext "enabled" "defaultSeccompProfile" }}
  {{ $key }}: {{- toYaml $value | nindent 4 }}
{{- end -}}
{{- if not .Values.podSecurityContext.seccompProfile -}}
{{- if .Values.apiGatewayNginx.podSecurityContext.defaultSeccompProfile -}}
{{- include "keptn.common.security-context-seccomp" . -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- else -}}
securityContext:
  fsGroup: 65532
{{- include "keptn.common.security-context-seccomp" . }}
{{- end -}}
{{- end -}}

{{- define "keptn.common.container-security-context" -}}
{{- if .Values.containerSecurityContext -}}
{{- if .Values.containerSecurityContext.enabled -}}
securityContext:
{{- range $key, $value := omit .Values.containerSecurityContext "enabled" }}
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
{{- define "keptn.common.update-strategy" -}}
{{- if .Values.strategy -}}
strategy:
{{- toYaml .Values.strategy | nindent 2 -}}
{{- else -}}
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1
    maxUnavailable: 0
{{- end -}}
{{- end -}}

{{/*
Renders nodeSelector if either "value" or "default" is not empty. the needed indentation must be set
with "indent". Templates may be used in values and keys!
Usage:
{{ include "keptn.nodeSelector" ( dict "value" .my-path-to.nodeselector-map "default" .my-path.to-default-nodeselector "indent" 6 "context" $ ) }}
*/}}
{{- define "keptn.nodeSelector" -}}
  {{- if not .indent }}
    {{ fail "keptn.nodeSelector needs indent to be set" }}
  {{- end }}
  {{- if not (typeIs "int" .indent) }}
    {{ fail "keptn.nodeSelector needs indent to be an int" }}
  {{- end }}
  {{- if or .value .default }}
    {{- printf "\n%snodeSelector:" (repeat .indent " ") }}{{- include "keptn.tpl-value-or-default" ( dict "value" .value "default" .default "context" .context  ) | nindent ( int ( add .indent 2 ) ) }}
  {{- end }}
{{- end -}}

{{/*
Renders affinity if either "value" or "default" is not empty. the needed indentation must be set
with "indent". Templates may be used in values and keys!
Usage:
{{ include "keptn.affinity" ( dict "value" .my-path-to.affinity-map "default" .my-path.to-default-affinity "indent" 6 "context" $ ) }}
*/}}
{{- define "keptn.affinity" -}}
  {{- if not .indent }}
    {{ fail "keptn.affinity needs indent to be set" }}
  {{- end }}
  {{- if not (typeIs "int" .indent) }}
    {{ fail "keptn.affinity needs indent to be an int" }}
  {{- end }}
  {{- if or .value .default }}
    {{- printf "\n%saffinity:" (repeat .indent " ") }}{{- include "keptn.tpl-value-or-default" ( dict "value" .value "default" .default "context" .context  ) | nindent ( int ( add .indent 2 ) ) }}
  {{- end }}
{{- end -}}

{{/*
Return a soft nodeAffinity definition 
{{ include "keptn.affinities.nodes.soft" (dict "preset" "FOO" "context" .context) -}}
*/}}
{{- define "keptn.affinities.nodes.soft" -}}
nodeAffinity:
  preferredDuringSchedulingIgnoredDuringExecution:
    - preference:
        matchExpressions:
          - key: {{ .preset.key }}
            operator: In
            values:
              {{- range .preset.values }}
              - {{ . }}
              {{- end }}
      weight: 1
{{- end -}}

{{/*
Return a hard nodeAffinity definition
{{ include "keptn.affinities.nodes.hard" (dict "preset" "FOO" "context" .context) -}}
*/}}
{{- define "keptn.affinities.nodes.hard" -}}
nodeAffinity:
  requiredDuringSchedulingIgnoredDuringExecution:
    nodeSelectorTerms:
      - matchExpressions:
          - key: {{ .preset.key }}
            operator: In
            values:
              {{- range .preset.values }}
              - {{ . }}
              {{- end }}
{{- end -}}

{{/*
Return a nodeAffinity definition
{{ include "keptn.affinities.nodes" (dict "value" "service-values" "default" "default-values" "component" "component-name" "context" . ) -}}
*/}}
{{- define "keptn.affinities.nodes" -}}
  {{- $preset := default "" .default -}}
  {{- if .value }}
    {{- $preset = .value -}}
  {{- end }}
  {{- if eq $preset.type "soft" }}
{{ include "keptn.affinities.nodes.soft" ( dict "preset" $preset "context" .context ) }}
  {{- else if eq $preset.type "hard" }}
{{ include "keptn.affinities.nodes.hard" ( dict "preset" $preset "context" .context ) }}
  {{- end }}
{{- end -}}

{{/*
Return a soft podAffinity/podAntiAffinity definition
{{ include "keptn.affinities.pods.soft" "component" "component-name" "context" . -}}
*/}}
{{- define "keptn.affinities.pods.soft" -}}
{{- $component := default "" .component -}}
preferredDuringSchedulingIgnoredDuringExecution:
  - podAffinityTerm:
      labelSelector:
        matchLabels: {{- (include "keptn.common.labels.selectorLabels" .context) | nindent 10 }}
          {{- if not (empty $component) }}
          app.kubernetes.io/name: {{ $component | quote }}
          {{- end }}
      namespaces:
        - {{ .context.Release.Namespace }}
      topologyKey: kubernetes.io/hostname
    weight: 1
{{- end -}}

{{/*
Return a hard podAffinity/podAntiAffinity definition
{{ include "keptn.affinities.pods.hard" "component" "component-name" "context" . -}}
*/}}
{{- define "keptn.affinities.pods.hard" -}}
{{- $component := default "" .component -}}
requiredDuringSchedulingIgnoredDuringExecution:
  - labelSelector:
      matchLabels: {{- (include "keptn.common.labels.selectorLabels" .context) | nindent 8 }}
        {{- if not (empty $component) }}
        app.kubernetes.io/name: {{ $component | quote }}
        {{- end }}
    namespaces:
      - {{ .context.Release.Namespace }}
    topologyKey: kubernetes.io/hostname
{{- end -}}

{{/*
Return a podAffinity/podAntiAffinity definition
{{ include "keptn.affinities.pods" (dict "value" "service-values" "default" "default-values" "component" "component-name" "context" . ) -}}
*/}}
{{- define "keptn.affinities.pods" -}}
  {{- $value := default "" .default -}}
  {{- if or .value.podAffinityPreset .value.podAntiAffinityPreset }}
    {{- $value = .value -}}
  {{- end }}
  {{- if and $value.podAffinityPreset ( not $value.podAntiAffinityPreset ) }}
podAffinity: {{- include "keptn.affinities.pods.mode" ( dict "mode" $value.podAffinityPreset "component" .component "context" .context ) | nindent 2 }}
  {{- else if and $value.podAntiAffinityPreset ( not $value.podAffinityPreset ) }}
podAntiAffinity: {{- include "keptn.affinities.pods.mode" ( dict "mode" $value.podAntiAffinityPreset "component" .component "context" .context ) | nindent 2 }}
  {{- end }}
{{- end -}}

{{/*
Return a soft or hard affinity definition
{{ include "keptn.affinities.pods.mode" (dict "mode" "soft or hard" "component" "component-name" "context" . ) -}}
*/}}
{{- define "keptn.affinities.pods.mode" -}}
{{- if .mode }}
  {{- if eq .mode "soft" }}
  {{- include "keptn.affinities.pods.soft" ( dict "component" .component "context"  .context ) }}
  {{- else if eq .mode "hard" }}
  {{- include "keptn.affinities.pods.hard" ( dict "component" .component "context"  .context ) }}
  {{- end }}
{{- end }}
{{- end -}}

{{/*
Renders tolerations if either "value" or "default" is not empty. the needed indentation must be set
with "indent". Templates may be used in values and keys!
Usage:
{{ include "keptn.tolerations" ( dict "value" .my-path-to.tolerations-map "default" .my-path.to-default-tolerations "indent" 6 "context" $ ) }}
*/}}
{{- define "keptn.tolerations" -}}
  {{- if not .indent }}
    {{ fail "keptn.tolerations needs indent to be set" }}
  {{- end }}
  {{- if not (typeIs "int" .indent) }}
    {{ fail "keptn.tolerations needs indent to be an int" }}
  {{- end }}
  {{- if or .value .default }}
    {{- printf "\n%stolerations:" (repeat .indent " ") }}{{- include "keptn.tpl-value-or-default" ( dict "value" .value "default" .default "context" .context  ) | nindent ( int ( add .indent 2 ) ) }}
  {{- end }}
{{- end -}}

{{/*
Renders a optional value that contains a template. if the given value is empty default is used
Usage:
{{ include "keptn.tpl-value-or-default" ( dict "value" .my-value.to-template "default" .my-default.to-template "context" $ ) }}
*/}}
{{- define "keptn.tpl-value-or-default" -}}
  {{- if .value }}
    {{- include "keptn.common.tplvalues.render" ( dict "value" .value "context" .context ) }}
  {{- else }}
    {{- include "keptn.common.tplvalues.render" ( dict "value" .default "context" .context ) }}
  {{- end }}
{{- end -}}

{{- define "keptn.initContainers.wait-for-nats" -}}
- name: "wait-for-nats"
  image: "{{ .Values.global.initContainers.image }}:{{ .Values.global.initContainers.tag }}"
  env:
  - name: "ENDPOINT"
    value: "http://{{ .Values.nats.nameOverride }}:8222"
  command: ['sh', '-c', 'until curl -s $ENDPOINT; do echo waiting for $ENDPOINT; sleep 2; done;']
  resources:
    limits:
      cpu: "50m"
      memory: "16Mi"
    requests:
      cpu: "25m"
      memory: "8Mi"
  securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    runAsUser: 65534
    capabilities:
      drop: ["ALL"]
{{- end -}}

{{- define "keptn.initContainers.wait-for-keptn-mongo" -}}
- name: "wait-for-keptn-mongo"
  image: "{{ .Values.global.initContainers.image }}:{{ .Values.global.initContainers.tag }}"
  env:
  - name: "ENDPOINT"
    value: "http://{{ .Release.Name }}-mongo:{{ .Values.mongo.service.ports.mongodb }}"
  command: ['sh', '-c', 'until curl -s $ENDPOINT; do echo waiting for $ENDPOINT; sleep 2; done;']
  resources:
    limits:
      cpu: "50m"
      memory: "16Mi"
    requests:
      cpu: "25m"
      memory: "8Mi"
  securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    runAsUser: 65534
    capabilities:
      drop: ["ALL"]
{{- end -}}

{{- define "keptn.initContainers.wait-for-mongodb-datastore" -}}
- name: "wait-for-mongodb-datastore"
  image: "{{ .Values.global.initContainers.image }}:{{ .Values.global.initContainers.tag }}"
  env:
  - name: "ENDPOINT"
    value: "http://mongodb-datastore:8080/health"
  command: ['sh', '-c', 'until curl -s $ENDPOINT; do echo waiting for $ENDPOINT; sleep 2; done;']
  resources:
    limits:
      cpu: "50m"
      memory: "16Mi"
    requests:
      cpu: "25m"
      memory: "8Mi"
  securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    runAsUser: 65534
    capabilities:
      drop: ["ALL"]
{{- end -}}

{{- define "keptn.initContainers.wait-for-shipyard-controller" -}}
- name: "wait-for-shipyard-controller"
  image: "{{ .Values.global.initContainers.image }}:{{ .Values.global.initContainers.tag }}"
  env:
  - name: "ENDPOINT"
    value: "http://shipyard-controller:8080/health"
  command: ['sh', '-c', 'until curl -s $ENDPOINT; do echo waiting for $ENDPOINT; sleep 2; done;']
  resources:
    limits:
      cpu: "50m"
      memory: "16Mi"
    requests:
      cpu: "25m"
      memory: "8Mi"
  securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    runAsUser: 65534
    capabilities:
      drop: ["ALL"]
{{- end -}}

{{- define "keptn.initContainers.wait-for-secret-service" -}}
- name: "wait-for-secret-service"
  image: "{{ .Values.global.initContainers.image }}:{{ .Values.global.initContainers.tag }}"
  env:
  - name: "ENDPOINT"
    value: "http://secret-service:8080/v1/secret"
  command: ['sh', '-c', 'until curl -s $ENDPOINT; do echo waiting for $ENDPOINT; sleep 2; done;']
  resources:
    limits:
      cpu: "50m"
      memory: "16Mi"
    requests:
      cpu: "25m"
      memory: "8Mi"
  securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    runAsUser: 65534
    capabilities:
      drop: ["ALL"]
{{- end -}}

{{- define "keptn.initContainers.wait-for-api-service" -}}
- name: "wait-for-secret-service"
  image: "{{ .Values.global.initContainers.image }}:{{ .Values.global.initContainers.tag }}"
  env:
  - name: "ENDPOINT"
    value: "http://api-service:8080/health"
  command: ['sh', '-c', 'until curl -s $ENDPOINT; do echo waiting for $ENDPOINT; sleep 2; done;']
  resources:
    limits:
      cpu: "50m"
      memory: "16Mi"
    requests:
      cpu: "25m"
      memory: "8Mi"
  securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    runAsUser: 65534
    capabilities:
      drop: ["ALL"]
{{- end -}}

{{- define "keptn.initContainers.wait-for-resource-service" -}}
- name: "wait-for-resource-service"
  image: "{{ .Values.global.initContainers.image }}:{{ .Values.global.initContainers.tag }}"
  env:
  - name: "ENDPOINT"
    value: "http://resource-service:8080/health"
  command: ['sh', '-c', 'until curl -s $ENDPOINT; do echo waiting for $ENDPOINT; sleep 2; done;']
  resources:
    limits:
      cpu: "50m"
      memory: "16Mi"
    requests:
      cpu: "25m"
      memory: "8Mi"
  securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    runAsUser: 65534
    capabilities:
      drop: ["ALL"]
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "keptn.imagePullSecrets" -}}
{{- include "keptn.common.images.renderPullSecrets" (dict "images" (list .Values.mongo.image .Values.shipyardController.image .Values.resourceService.image .Values.approvalService.image .Values.lighthouseService.image .Values.mongodbDatastore.image .Values.remediationService.image .Values.secretService.image .Values.statisticsService.image .Values.webhookService.image .Values.apiGatewayNginx.image .Values.apiService.image .Values.bridge.image .Values.distributor.image) "context" $) -}}
{{- end -}}

{{- define "keptn.mongodb-credentials.volume" -}}
- name: mongodb-credentials
  secret:
    secretName: mongodb-credentials
    defaultMode: 0400
    items:
      - key: mongodb-user
        path: mongodb-user
      - key: mongodb-passwords
        path: mongodb-passwords
      - key: external_connection_string
        path: external_connection_string
{{- end -}}

{{- define "keptn.mongodb-credentials.volumeMount" -}}
- name: mongodb-credentials
  mountPath: /config/mongodb_credentials
  readOnly: true
{{- end -}}