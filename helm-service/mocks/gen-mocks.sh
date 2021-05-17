#!/bin/bash

GO111MODULE=off go get -u github.com/golang/mock/mockgen

set -euo pipefail
set -x

mockgen -package mocks -destination=./mock_chart_storer.go github.com/keptn/keptn/helm-service/pkg/types IChartStorer
mockgen -package mocks -destination=./mock_chart_packager.go github.com/keptn/keptn/helm-service/pkg/types IChartPackager
mockgen -package mocks -destination=./mock_namespace_manager.go github.com/keptn/keptn/helm-service/pkg/namespacemanager INamespaceManager
mockgen -package mocks -destination=./mock_chart_generator.go github.com/keptn/keptn/helm-service/pkg/helm ChartGenerator
mockgen -package mocks -destination=./mock_configuration_changer.go github.com/keptn/keptn/helm-service/pkg/configurationchanger IConfigurationChanger
mockgen -package mocks -destination=./mock_helm_executor.go github.com/keptn/keptn/helm-service/pkg/helm HelmExecutor
mockgen -package mocks -destination=./mock_project_operator.go github.com/keptn/keptn/helm-service/pkg/types IProjectHandler
mockgen -package mocks -destination=./mock_stages_handler.go  github.com/keptn/keptn/helm-service/pkg/types IStagesHandler
mockgen -package mocks -destination=./mock_mesh.go github.com/keptn/keptn/helm-service/pkg/mesh Mesh
mockgen -package mocks -destination=./mock_service_handler.go github.com/keptn/keptn/helm-service/pkg/types IServiceHandler
mockgen -package mocks -destination=./mock_resource_handler.go github.com/keptn/keptn/helm-service/pkg/types IResourceHandler
mockgen -package mocks -destination=./mock_onboarder.go github.com/keptn/keptn/helm-service/controller Onboarder