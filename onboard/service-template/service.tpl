apiVersion: v1
kind: Service
metadata:
  name: {{ .Chart.Name }}-SERVICE_PLACEHOLDER_DEC
  labels:
    chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
spec:
  type: {{ .Values.SERVICE_PLACEHOLDER_C.service.type }}
  ports:
  - port: {{ .Values.SERVICE_PLACEHOLDER_C.service.externalPort }}
    targetPort: {{ .Values.SERVICE_PLACEHOLDER_C.service.internalPort }}
    protocol: TCP
    name: {{ .Values.SERVICE_PLACEHOLDER_C.service.name }}
  selector:
    app: {{ .Chart.Name }}-selector-SERVICE_PLACEHOLDER_DEC