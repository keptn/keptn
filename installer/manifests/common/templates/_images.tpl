{{/*Return the proper serivce image name*/}}
{{/*{{- include "keptn.common.images.image" ( dict "imageRoot" .Values.helmservice.image "global" .Values.global.keptn "defaultTag" .Chart.appVersion) -}}*/}}
{{- define "keptn.common.images.image" -}}
{{- $registryName := "" -}}
{{- $repositoryName := .imageRoot.repository -}}
{{- $tag := include "keptn.common.images.tag" (dict "imageRoot" .imageRoot "global" .global "defaultTag" .defaultTag) -}}
{{- if .global }}
    {{- if .global.registry }}
     {{- $registryName = .global.registry -}}
    {{- end -}}
{{- end -}}
{{- if .imageRoot.registry -}}
  {{- $registryName = .imageRoot.registry -}}
{{- end -}}
{{- if $registryName }}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- else -}}
{{- printf "%s:%s" $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*Return the proper serivce image tag*/}}
{{/*{{- include "keptn.common.images.tag" ( dict "imageRoot" .Values.helmservice.image "global" .Values.global.keptn "defaultTag" .Chart.appVersion) -}}*/}}
{{- define "keptn.common.images.tag" -}}
{{- $tag := "" -}}
{{/* Set Image Tag: if globally set or at service level or use default from Chart.yaml*/}}
{{- if .global -}}
    {{- if .global.tag -}}
       {{- $tag = .global.tag -}}
    {{- end -}}
{{- end -}}
{{- if .imageRoot.tag -}}
   {{- $tag = .imageRoot.tag -}}
{{- end -}}
{{- if not $tag }}
   {{- $tag = .defaultTag -}}
{{- end -}}
{{- printf "%s" $tag -}}
{{- end -}}
