{{/* Expand the name of the chart. */}}
{{- define "syphon.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "syphon.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/* Create chart name and version as used by the chart label. */}}
{{- define "syphon.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Common labels */}}
{{- define "syphon.labels" -}}
helm.sh/chart: {{ include "syphon.chart" . }}
{{ include "syphon.selectorLabels" . }}
app.kubernetes.io/version: {{ .image.tag | default .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/* Selector labels */}}
{{- define "syphon.selectorLabels" -}}
app.kubernetes.io/name: {{ include "syphon.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}


{{/* Define name for configmap that changes on content change */}}
{{- define "syphon.configMapName" -}}
name: {{ $.Release.Name }}-{{ include "syphon.config" $ | sha256sum | trunc 10 }}
{{- end -}}

{{/* Convert config to yaml */}}
{{- define "syphon.config" -}}
config.yaml: |-
  {{- toToml .Values.config | trim | nindent 2 }}
{{- end -}}
