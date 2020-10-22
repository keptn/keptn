#!/bin/bash

GO111MODULE=off go get -u github.com/golang/mock/mockgen
EXT_TYPES=github.com/keptn/keptn/helm-service/pkg/types

set -euo pipefail
set -x

mockgen -package mocks -destination=./mock_chart_storer.go github.com/keptn/kubernetes-utils/pkg ChartStorer
mockgen -package mocks -destination=./mock_chart_packager.go github.com/keptn/kubernetes-utils/pkg ChartPackager
mockgen -package mocks -destination=./mock_namespace_manager.go github.com/keptn/keptn/helm-service/pkg/namespacemanager INamespaceManager
mockgen -package mocks -destination=./mock_chart_generator.go github.com/keptn/keptn/helm-service/pkg/helm ChartGenerator
mockgen -package mocks -destination=./mock_configuration_changer.go github.com/keptn/keptn/helm-service/pkg/configurationchanger IConfigurationChanger

mockgen -package mocks -destination=./mock_project_operator.go ${EXT_TYPES} ProjectOperator
mockgen -package mocks -destination=./mock_stages_handler.go  ${EXT_TYPES} IStagesHandler
mockgen -package mocks -destination=./mock_mesh.go ${EXT_TYPES} Mesh
mockgen -package mocks -destination=./mock_service_handler.go ${EXT_TYPES} IServiceHandler