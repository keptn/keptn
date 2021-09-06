{{/*
Expand the name of the chart.
*/}}
{{- define "jmeter-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "jmeter-service.fullname" -}}
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
{{- define "jmeter-service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "jmeter-service.labels" -}}
helm.sh/chart: {{ include "jmeter-service.chart" . }}
{{ include "jmeter-service.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{/*
Selector labels
*/}}
{{- define "jmeter-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "jmeter-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "jmeter-service.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "jmeter-service.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Return the proper image name
{{ dict "imageRoot" .Values.path.to.the.image "context" $ | include "jmeter-service.images.image" }}
*/}}
{{- define "jmeter-service.images.image" -}}
{{- $global := .context.Values.global -}}
{{- $registryName := .imageRoot.registry -}}
{{- $repositoryName := .imageRoot.repository -}}
{{- $tag := default .context.Chart.AppVersion .imageRoot.tag | toString -}}
{{- if $global }}
    {{- if $global.imageRegistry }}
     {{- $registryName = $global.imageRegistry -}}
    {{- end -}}
{{- end -}}
{{- if $registryName }}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- else -}}
{{- printf "%s:%s" $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names evaluating values as templates
{{ list .Values.path.to.the.image1 .Values.path.to.the.image2 | dict "context" $ "images" | include "jmeter-service.images.renderPullSecrets" }}
*/}}
{{- define "jmeter-service.images.renderPullSecrets" -}}
  {{- $pullSecrets := list }}
  {{- $context := .context }}

  {{- if $context.Values.global }}
    {{- range $context.Values.global.imagePullSecrets -}}
      {{- $pullSecrets = dict "value" . "context" $context | include "jmeter-service.tplvalues.render" | append $pullSecrets -}}
    {{- end -}}
  {{- end -}}

  {{- range .images -}}
    {{- range .pullSecrets -}}
      {{- $pullSecrets = dict "value" . "context" $context | include "jmeter-service.tplvalues.render" | append $pullSecrets -}}
    {{- end -}}
  {{- end -}}

  {{- if empty $pullSecrets | not -}}
imagePullSecrets:
    {{- range $pullSecrets }}
- name: {{ . }}
    {{- end }}
  {{- end }}
{{- end }}

{{/*
Renders a value that contains template.
Usage:
{{ dict "value" .Values.path.to.the.Value "context" $ | include "jmeter-service.tplvalues.render" }}
*/}}
{{- define "jmeter-service.tplvalues.render" -}}
    {{- if typeIs "string" .value }}
        {{- tpl .value .context }}
    {{- else }}
        {{- tpl (toYaml .value) .context }}
    {{- end }}
{{- end }}

{{/*
Return the image name.
Usage:
{{ dict "space" .Values.imageNameSpace "context" $ | include "jmeter-service.image.name" }}
*/}}
{{- define "jmeter-service.image.name" }}
{{- dict "imageRoot" .space.image "context" .context | include "jmeter-service.images.image" }}
{{- end }}

{{/*
Return the proper Docker Image Registry Secret Names.
Usage:
{{ list .Values.imageNameSpace0 .Values.imageNameSpace1 | dict "context" $ "indent" $number "bases" | "jmeter-service.image.pullSecrets" }}
*/}}
{{- define "jmeter-service.image.pullSecrets" }}
{{- $images := list }}
{{- range .bases }}
  {{- $images = append $images .image }}
{{- end }}
{{- $content := dict "images" $images "context" .context | include "jmeter-service.images.renderPullSecrets" }}
{{- if $content }}
  {{- nindent .indent $content }}
{{- end }}
{{- end }}
