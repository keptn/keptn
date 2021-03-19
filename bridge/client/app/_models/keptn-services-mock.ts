import {KeptnService} from './keptn-service';
import {Subscription} from './subscription';

const services: KeptnService[] = [
  {
    name: 'ansible-service',
    version: '0.2.0',
    host: 'gke_research_us-central1-c_keptn-exec',
    namespace: 'keptn-uniform',
    location: 'Execution plane-B',
    status: 'healthy',
    subscriptions: []
  } as KeptnService,
  {
    name: 'helm-service',
    version: '0.8.0',
    host: 'gke_research_us-central1-c_keptn-exec',
    namespace: 'keptn-uniform',
    location: 'Execution plane-A',
    status: 'healthy',
    subscriptions: []
  } as KeptnService,
  {
    name: 'jenkins-service',
    version: '0.8.0',
    host: 'gke_research_us-central1-c_keptn-exec',
    namespace: 'keptn-uniform',
    location: 'Execution plane-A',
    status: 'healthy',
    subscriptions: []
  } as KeptnService,
  {
    name: 'dynatrace-sli-service',
    version: '0.7.3',
    host: 'gke_research_us-central1-c_keptn-control',
    namespace: 'keptn',
    location: 'Control plane',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Fetch SLI',
        event: 'get-sli.triggered',
        stages: ['development', 'staging'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'secret', value: 'secret_name'},
          {key: 'pipeline', value: 'deploy_service_pipeline'},
        ]
      } as Subscription,
      {
        name: 'Deployment in production',
        event: 'deployment.triggered',
        stages: ['production'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'secret2', value: 'secret_name2'},
          {key: 'pipeline2', value: 'deploy_service_pipeline2'},
        ]
      } as Subscription
    ]
  } as KeptnService,
  {
    name: 'dynatrace-service',
    version: '0.11.0',
    host: 'gke_research_us-central1-c_keptn-control',
    namespace: 'keptn',
    location: 'Control plane',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Deployment in staging',
        event: 'deployment.triggered',
        stages: ['staging'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'secret', value: 'secret_name'},
          {key: 'pipeline', value: 'deploy_service_pipeline'}
        ]
      } as Subscription,
      {
        name: 'Deployment in production',
        event: 'deployment.triggered',
        stages: ['production'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'secret', value: 'secret_name'},
          {key: 'pipeline', value: 'deploy_service_pipeline'}
        ]
      } as Subscription,
      {
        name: 'Test in staging',
        event: 'test.triggered',
        stages: ['staging'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'secret', value: 'secret_name'},
          {key: 'pipeline', value: 'deploy_service_pipeline'}
        ]
      } as Subscription,
      {
        name: 'Evaluation',
        event: 'evaluation.triggered',
        stages: ['development', 'staging'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'dynatraceSecret', value: 'xyz12345'}
        ]
      } as Subscription,
      {
        name: 'All releases',
        event: 'release.triggered',
        stages: ['development', 'staging', 'production'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'dynatraceSecret', value: 'xyz12345'}
        ]
      } as Subscription,
      {
        name: 'All rollbacks',
        event: 'rollback.triggered',
        stages: ['development', 'staging', 'production'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'dynatraceSecret', value: 'xyz12345'}
        ]
      } as Subscription,
    ]
  } as KeptnService,
  {
    name: 'servicenow-service',
    version: '0.2.0',
    host: 'gke_research_us-central1-c_keptn-control',
    namespace: 'keptn',
    location: 'Control plane',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Fetch SLI',
        event: 'get-sli.triggered',
        stages: ['development', 'staging'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'dynatraceSecret', value: 'xyz12345'}
        ]
      } as Subscription,
      {
        name: 'Deployment in development',
        event: 'deployment.triggered',
        stages: ['development'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'dynatraceSecret', value: 'abcde12345'}
        ]
      } as Subscription
    ]
  } as KeptnService
];

const KeptnServicesMock = JSON.parse(JSON.stringify(services));
export {KeptnServicesMock};
