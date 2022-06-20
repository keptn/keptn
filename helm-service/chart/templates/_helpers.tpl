{{/*
Expand the name of the chart.
*/}}
{{- define "helm-service.name" -}}
{{- include "common.names.name" . -}}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "helm-service.fullname" -}}
{{- include "common.names.fullname" . -}}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "helm-service.chart" -}}
{{- include "common.names.chart" . -}}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "helm-service.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "helm-service.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
preStop hook for helm-service deployment
*/}}
{{- define "helm-service.prestop" -}}
lifecycle:
  preStop:
    exec:
      # using 60s of sleeping to be on the safe side before terminating the pod
      command: ["/bin/sleep", {{ .Values.helmservice.preStopHookTime | default 60 | quote }} ]
{{- end }}

