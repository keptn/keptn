{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "keptn.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "keptn.fullname" -}}
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
{{- define "keptn.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "keptn.labels" -}}
helm.sh/chart: {{ include "keptn.chart" . }}
{{ include "keptn.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "keptn.selectorLabels" -}}
app.kubernetes.io/name: {{ include "keptn.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "keptn.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "keptn.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "keptn.distributor.resources" -}}
resources:
  requests:
    memory: "16Mi"
    cpu: "25m"
  limits:
    memory: "32Mi"
    cpu: "100m"
{{- end }}

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
Renders a optional value that contains a template. if the given value is empty default is used
Usage:
{{ include "keptn.tpl-value-or-default" ( dict "value" .my-value.to-template "default" .my-default.to-template "context" $ ) }}
*/}}
{{- define "keptn.tpl-value-or-default" -}}
  {{- if .value }}
    {{- include "keptn.tpl-value" ( dict "value" .value "context" .context ) }}
  {{- else }}
    {{- include "keptn.tpl-value" ( dict "value" .default "context" .context ) }}
  {{- end }}
{{- end -}}

{{/*
Renders a value that contains a template. value can be a string, map or array
Usage:
{{ include "keptn.tpl-value" ( dict "value" .my-value.to-template "context" $ ) }}
*/}}
{{- define "keptn.tpl-value" -}}
  {{- if typeIs "string" .value }}
    {{- tpl .value .context }}
  {{- else }}
    {{- tpl (.value | toYaml) .context }}
  {{- end }}
{{- end -}}
