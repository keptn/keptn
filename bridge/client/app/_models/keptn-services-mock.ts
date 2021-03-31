import {KeptnService} from './keptn-service';
import {Subscription} from './subscription';

const services: KeptnService[] = [
  {
    name: 'ansible-service',
    version: '0.2.0',
    host: 'gke_research_us-central1-c_prod-customer-A',
    namespace: 'keptn-uniform',
    location: 'Execution plane-A',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Remediation action',
        event: 'action.triggered',
        stages: [],
        services: [],
        parameters: [
          {key: 'ansibleTowerSecret', value: 'ansibleSecret'},
        ]
      } as Subscription
    ]
  } as KeptnService,
  {
    name: 'helm-service',
    version: '0.8.1',
    host: 'gke_research_us-central1-c_prod-customer-A',
    namespace: 'keptn-uniform',
    location: 'Execution plane-A',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Deployments in Production - Customer A',
        event: 'deployment.triggered',
        stages: [],
        services: [],
        parameters: [
          {key: 'helmSecret', value: 'helmSecret'},
        ]
      } as Subscription
    ]
  } as KeptnService,
  {
    name: 'jmeter-service',
    version: '0.8.1',
    host: 'gke_research_us-central1-c_prod-customer-A',
    namespace: 'keptn-uniform',
    location: 'Execution plane-A',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Tests in Production - Customer A',
        event: 'test.triggered',
        stages: ['development', 'staging'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'secret', value: 'secret_name'}
        ]
      } as Subscription
    ]
  } as KeptnService,
  {
    name: 'lighthouse-service',
    version: '0.8.0',
    host: 'gke_research_us-central1-c_keptn-control',
    namespace: 'keptn',
    location: 'Control plane',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Evaluation',
        event: 'evaluation.triggered',
        stages: [],
        services: [],
        parameters: [
          {key: 'lighthouseSecret', value: 'xyz12345'}
        ]
      } as Subscription
    ]
  } as KeptnService,
  {
    name: 'approval-service',
    version: '0.8.0',
    host: 'gke_research_us-central1-c_keptn-control',
    namespace: 'keptn',
    location: 'Control plane',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Approval events',
        event: 'approval.triggered',
        stages: [ ],
        services: [ ],
        parameters: [ ]
      } as Subscription
    ]
  } as KeptnService,
  {
    name: 'remediation-service',
    version: '0.8.0',
    host: 'gke_research_us-central1-c_keptn-control',
    namespace: 'keptn',
    location: 'Control plane',
    status: 'healthy',
    subscriptions: []
  } as KeptnService,
  {
    name: 'jenkins-service',
    version: '0.8.0',
    host: 'aws_us-central1-c_prod-customer-B',
    namespace: 'keptn-uniform',
    location: 'Execution plane-B',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Deployments of cards, cards-db on Production',
        event: 'deployment.triggered',
        stages: ['production'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'jenkinsPWD', value: 'pwd'},
          {key: 'pipeline', value: 'pipeline'}
        ]
      } as Subscription
    ]
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
        name: 'Get SLI',
        event: 'get-sli.triggered',
        stages: ['development', 'staging'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'dynatraceToken', value: 'token'}
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
        name: 'Listen to all Events',
        event: 'all',
        stages: ['staging'],
        services: ['cards', 'cards-db'],
        parameters: [
          {key: 'secret', value: 'secret_name'},
        ]
      } as Subscription
    ]
  } as KeptnService,
  {
    name: 'servicenow-service',
    version: '0.2.0',
    host: 'gke_research_us-central1-c_prod-customer-A',
    namespace: 'keptn-uniform',
    location: 'Execution plane-A',
    status: 'healthy',
    subscriptions: [
      {
        name: 'Remediation action',
        event: 'action.triggered',
        stages: [],
        services: [],
        parameters: [
          {key: 'snowSecret', value: 'secret'},
        ]
      } as Subscription
    ]
  } as KeptnService
];

const KeptnServicesMock = JSON.parse(JSON.stringify(services));
export {KeptnServicesMock};
