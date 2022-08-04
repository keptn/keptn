import { EvaluationFinishedMock } from '../fixtures/typed/evaluationFinished.mock';
import { EvaluationFinishedScoredMock } from '../fixtures/typed/evaluationFinishedScoreMock';

export function interceptEmptyEnvironmentScreen(): void {
  interceptProjectBoard();
  cy.intercept('/api/project/dynatrace?approval=true&remediation=true', { fixture: 'project.empty.mock' }).as(
    'project'
  );
  cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
    fixture: 'get.projects.empty.mock',
  }).as('projects');
  cy.intercept('GET', '/api/project/dynatrace/services', {
    statusCode: 200,
    body: [],
  });
  cy.intercept('POST', '/api/controlPlane/v1/project/dynatrace/service', {
    statusCode: 200,
    body: {},
  });
}

export function interceptEnvironmentScreen(): void {
  const project = 'sockshop';
  interceptProjectBoard();
  cy.intercept('/api/project/sockshop?approval=true&remediation=true', { fixture: 'project.mock' }).as('project');
  cy.intercept('/api/project/sockshop/customSequences', { fixture: 'custom-sequences.mock' }).as('customSequences');
  cy.intercept('POST', '/api/v1/event', { body: { keptnContext: '6c98fbb0-4c40-4bff-ba9f-b20556a57c8a' } });
  cy.intercept('POST', '/api/controlPlane/v1/project/sockshop/stage/dev/service/carts/evaluation', {
    body: { keptnContext: '6c98fbb0-4c40-4bff-ba9f-b20556a57c8a' },
  });

  cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=1&keptnContext=6c98fbb0-4c40-4bff-ba9f-b20556a57c8a', {
    fixture: 'eventByContext.mock',
  }).as('triggeredSequence');
  cy.intercept(
    '/api/controlPlane/v1/sequence/sockshop?pageSize=10&fromTime=2022-02-23T14:28:50.504Z&beforeTime=2021-07-06T09:22:56.433Z',
    {
      body: {
        states: [],
      },
    }
  ).as('betweenTriggeredSequence');
  cy.intercept('/api/mongodb-datastore/event?keptnContext=6c98fbb0-4c40-4bff-ba9f-b20556a57c8a&project=sockshop*', {
    body: {
      events: [],
      pageSize: 100,
      totalCount: 0,
    },
  }).as('triggeredSequenceEvents');

  setEvaluationUrls(project, 'carts');
  setEvaluationUrls(project, 'carts-db');
}

function setEvaluationUrls(project: string, service: string): void {
  for (const stage of ['dev', 'staging', 'production']) {
    for (const limit of [5, 6]) {
      cy.intercept(
        'GET',
        `/api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:${project}%20AND%20data.service:${service}%20AND%20data.stage:${stage}%20AND%20source:lighthouse-service&excludeInvalidated=true&limit=${limit}`,
        { body: { events: [] } }
      ).as(`evaluationHistory-${service}-${stage}-${limit}`);
    }
  }
}

export function interceptMainResourceEnabled(): void {
  cy.intercept('/api/v1/metadata', { fixture: 'metadata.ap-disabled.mock' }).as('metadata');
  cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoEnableResourceService.mock' });
  cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' }).as(
    'projects'
  );
}

export function interceptMainResourceApEnabled(): void {
  cy.intercept('/api/v1/metadata', { fixture: 'metadata.ap-enabled.mock' }).as('metadata');
  cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoEnableResourceService.mock' });
  cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' }).as(
    'projects'
  );
}

export function interceptMain(): void {
  cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' }).as('metadata');
  cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
  cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' }).as(
    'projects'
  );
}

export function interceptFailedMetadata(): void {
  cy.intercept('/api/v1/metadata', { forceNetworkError: true }).as('metadata');
}

export function interceptCreateProject(): void {
  cy.intercept('POST', 'api/controlPlane/v1/project', {
    statusCode: 200,
    body: {},
  });
}

export function interceptProjectSettings(): void {
  cy.intercept('PUT', 'api/controlPlane/v1/project', {
    statusCode: 200,
    body: {},
  });
  cy.intercept('/api/project/sockshop', { fixture: 'project.mock' }).as('projectPlain');
  cy.intercept('DELETE', '/api/controlPlane/v1/project/sockshop', {
    statusCode: 200,
  });
}

export function interceptDashboard(): void {
  interceptMain();
  cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=5', { fixture: 'sequences.sockshop' }).as('sequences');
  cy.intercept('/api/controlPlane/v1/sequence/my-error-project?pageSize=5', { body: { states: [] } });
  cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' }).as(
    'projects'
  );
}

export function interceptProjectBoard(): void {
  interceptMain();
  cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: false });
}

export function interceptServicesPage(): void {
  cy.intercept('GET', '/api/project/sockshop/serviceStates', {
    statusCode: 200,
    fixture: 'get.sockshop.service.states.mock.json',
  }).as('serviceStates');
  cy.intercept('GET', '/api/project/sockshop/deployment/da740469-9920-4e0c-b304-0fd4b18d17c2', {
    statusCode: 200,
    fixture: 'get.sockshop.service.carts.deployment.mock.json',
  }).as('ServiceDeployment');
  cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
    statusCode: 200,
    fixture: 'get.sockshop.service.carts.evaluations.mock.json',
  }).as('serviceDatastore');
}

export function interceptServicesPageWithLoadingSequences(): void {
  interceptServicesPage();
  cy.intercept('GET', '/api/project/sockshop/deployment/da740469-9920-4e0c-b304-0fd4b18d17c2', {
    statusCode: 200,
    fixture: 'get.sockshop.service.carts.deployment.running.mock.json',
  }).as('ServiceDeployment');
}

export function interceptServicesPageWithRemediation(): void {
  interceptServicesPage();
  cy.intercept('GET', '/api/project/sockshop/serviceStates', {
    statusCode: 200,
    fixture: 'get.sockshop.service.states.with.remediation.mock.json',
  }).as('serviceStates');

  cy.intercept('GET', '/api/project/sockshop/deployment/fa66eea5-53a8-45b6-aefe-ef03c08b61e4', {
    statusCode: 200,
    fixture: 'get.sockshop.service.carts.deployment.with.remediations.mock.json',
  }).as('ServiceDeployment');
}

export function interceptSequencesPage(): void {
  interceptProjectBoard();
  cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', { fixture: 'sequences.sockshop' }).as('Sequences');
  cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=10&beforeTime=2021-07-06T09:22:56.433Z', {
    fixture: 'sequences-page-2.sockshop',
  }).as('SequencesPage2');
  cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=10&beforeTime=2021-07-06T08:13:53.766Z', {
    fixture: 'sequences-page-3.sockshop',
  }).as('SequencesPage3');
  cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25&fromTime=*', {
    body: {
      states: [],
    },
  }).as('SequencesUpdate');

  cy.intercept('/api/project/sockshop/sequences/filter', { fixture: 'sequence.filter.mock' }).as('SequencesMetadata');
  cy.intercept('/api/mongodb-datastore/event?keptnContext=62cca6f3-dc54-4df6-a04c-6ffc894a4b5e&project=sockshop', {
    fixture: 'sequence.traces.mock.json',
  });

  cy.intercept('/api/mongodb-datastore/event?keptnContext=99a20ef4-d822-4185-bbee-0d7a364c213b&project=sockshop', {
    fixture: 'sequence-traces/approval.mock.json',
  });

  cy.intercept('/api/controlPlane/v1/project/sockshop/stage/production/service/carts', {
    deployedImage: 'myImage:0.0.1',
  }).as('approvalImage');

  cy.intercept(
    '/api/mongodb-datastore/event?keptnContext=62cca6f3-dc54-4df6-a04c-6ffc894a4b5e&project=sockshop&fromTime=*',
    {
      body: [],
    }
  );
  interceptEvaluationOfApproval();
}

export function interceptEvaluationOfApproval(includeData = false, delay = 0): void {
  const data = includeData
    ? { fixture: 'get.approval-evaluation.mock.json' }
    : {
        body: {
          events: [],
        },
      };
  cy.intercept(
    '/api/mongodb-datastore/event?keptnContext=99a20ef4-d822-4185-bbee-0d7a364c213b&type=sh.keptn.event.evaluation.finished&source=lighthouse-service&stage=production&pageSize=1',
    {
      ...data,
      delay,
    }
  );
}

export function interceptSequencesPageWithSequenceThatIsNotLoaded(): void {
  interceptSequencesPage();
  const keptnContext = '1663de8a-a414-47ba-9566-10a9730f40ff';
  cy.intercept(`/api/mongodb-datastore/event?keptnContext=${keptnContext}&project=sockshop`, {
    fixture: 'sequence.traces.mock.json',
  }).as('sequenceTraces');

  cy.intercept(`/api/controlPlane/v1/sequence/sockshop?pageSize=1&keptnContext=${keptnContext}`, {
    fixture: 'get.sequence.mock.json',
  });

  cy.intercept(
    '/api/controlPlane/v1/sequence/sockshop?pageSize=10&fromTime=2021-07-06T08:13:53.766Z&beforeTime=2021-07-06T09:22:56.433Z',
    {
      fixture: 'get.sequence.mock.json',
    }
  );

  cy.intercept(`/api/mongodb-datastore/event?keptnContext=${keptnContext}&project=sockshop&fromTime=*`, {
    body: [],
  });
}

export function interceptIntegrations(): void {
  interceptProjectBoard();
  cy.intercept('/api/uniform/registration', { fixture: 'registration.mock' });
  // jmeter-service
  cy.intercept('/api/controlPlane/v1/log?integrationId=355311a7bec3f35bf3abc2484ab09bcba8e2b297&pageSize=100', {
    body: {
      logs: [],
    },
  });
  // approval-service
  cy.intercept('/api/controlPlane/v1/log?integrationId=4d57b2af3cdd66bce06625daafa9c5cbb474a6b8&pageSize=100', {
    body: {
      logs: [],
    },
  });
  // webhook-service
  cy.intercept('/api/controlPlane/v1/log?integrationId=0f2d35875bbaa72b972157260a7bd4af4f2826df&pageSize=100', {
    body: {
      logs: [
        {
          integrationid: '0f2d35875bbaa72b972157260a7bd4af4f2826df',
          message: 'my error',
          shkeptncontext: ' 7394b5b3-2fb3-4cb7-b435-d0e9d6f0cb87',
          task: 'my task',
          time: '2022-02-09T16:27:02.678Z',
          triggeredid: 'bd3bc477-6d0f-4d71-b15d-c33e953a74ba',
        },
      ],
    },
  });
  // jmeter-service uniform-info
  cy.intercept('/api/uniform/registration/355311a7bec3f35bf3abc2484ab09bcba8e2b297/info', {
    body: {
      isControlPlane: true,
      isWebhookService: false,
    },
  }).as('jmeterUniformInfo');
  cy.intercept('/api/uniform/registration/0f2d35875bbaa72b972157260a7bd4af4f2826df/info', {
    body: {
      isControlPlane: true,
      isWebhookService: true,
    },
  });
  cy.intercept('/api/uniform/registration/355311a7bec3f35bf3abc2484ab09bcba8e2b297/subscription', {
    body: {
      id: '0b77c90e-282d-4a7e-a96d-e23027265868',
    },
  });
  cy.intercept('/api/uniform/registration/0f2d35875bbaa72b972157260a7bd4af4f2826df/subscription', {
    body: {
      id: 'b5111b1c-446a-410d-bb6c-e1dcd409c890',
    },
  });
  cy.intercept('/api/uniform/registration/webhook-service/config', { body: true });
  cy.intercept('/api/project/sockshop/tasks', { fixture: 'tasks.mock' });
  cy.intercept('/api/secrets/scope/keptn-webhook-service', { fixture: 'secrets.mock' }).as('webhook-secrets');
  cy.intercept('/api/intersectEvents', { fixture: 'intersected-event.mock' });
  cy.intercept(
    'DELETE',
    '/api/uniform/registration/355311a7bec3f35bf3abc2484ab09bcba8e2b297/subscription/0e021b71-1533-4cfe-875a-b756aa6107ba?isWebhookService=false',
    {
      body: {},
    }
  );
  cy.intercept(
    '/api/controlPlane/v1/uniform/registration/355311a7bec3f35bf3abc2484ab09bcba8e2b297/subscription/0e021b71-1533-4cfe-875a-b756aa6107ba',
    { fixture: 'jmeter.mock' }
  );
  cy.intercept(
    '/api/uniform/registration/355311a7bec3f35bf3abc2484ab09bcba8e2b297/subscription/0e021b71-1533-4cfe-875a-b756aa6107ba',
    { body: {} }
  );

  cy.intercept('/api/project/sockshop/sequences/filter', { fixture: 'sequence.filter.mock' }).as('SequencesMetadata');
}

export function interceptNoWebhookSecrets(): void {
  cy.intercept('/api/secrets/scope/keptn-webhook-service', {
    body: [],
  }).as('webhook-secrets');
}

export function interceptSecrets(): void {
  interceptProjectBoard();

  cy.intercept('POST', 'api/secrets/v1/secret', {
    statusCode: 200,
  }).as('postSecrets');

  cy.intercept('GET', 'api/secrets/v1/secret', {
    statusCode: 200,
    body: {
      Secrets: [
        { name: 'dynatrace', scope: 'dynatrace-service', keys: ['DT_API_TOKEN', 'DT_TENANT'] },
        { name: 'dynatrace-prod', scope: 'dynatrace-service', keys: ['DT_API_TOKEN'] },
        { name: 'api', scope: 'keptn-default', keys: ['API_TOKEN'] },
        { name: 'webhook', scope: 'keptn-webhook-service', keys: ['webhook_url', 'webhook_secret', 'webhook_proxy'] },
      ],
    },
  }).as('getSecrets');

  cy.intercept('GET', 'api/project/sockshop', {
    statusCode: 200,
    fixture: 'get.approval.json',
  });

  cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
    statusCode: 200,
  }).as('hasUnreadUniformRegistrationLogs');

  cy.intercept('POST', 'api/uniform/registration', {
    statusCode: 200,
    body: '[]',
  }).as('uniformRegPost');

  cy.intercept('DELETE', 'api/secrets/v1/secret?name=dynatrace-prod&scope=dynatrace-service', {
    statusCode: 200,
  }).as('deleteSecret');

  cy.intercept('GET', 'api/secrets/v1/scope', {
    statusCode: 200,
    body: {
      scopes: ['dynatrace-service', 'keptn-webhook-service'],
    },
  });
}

export function interceptEvaluationBoardDynatrace(): void {
  cy.intercept('api/mongodb-datastore/event?keptnContext=*&type=sh.keptn.event.evaluation.triggered&pageSize=1', {
    fixture: 'service/get.evaluation.triggered.mock.json',
  });

  cy.intercept(
    'api/mongodb-datastore/event?keptnContext=*&type=sh.keptn.event.evaluation.finished&source=lighthouse-service',
    {
      fixture: 'service/get.event2.data.json',
    }
  );
  cy.intercept('api/controlPlane/v1/project/dynatrace/stage/quality-gate/service/items', {
    fixture: 'get.service.items.mock.json',
  });
}

export function interceptEvaluationBoard(): void {
  interceptMain();
  cy.intercept('api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
    fixture: 'service/get.eval.data.json',
  });

  cy.intercept('api/mongodb-datastore/event?keptnContext=*&type=sh.keptn.event.evaluation.triggered&pageSize=1', {
    fixture: 'service/get.evaluation.triggered-with-deployment.mock.json',
  });

  cy.intercept(
    'api/mongodb-datastore/event?keptnContext=*&type=sh.keptn.event.evaluation.finished&source=lighthouse-service',
    {
      fixture: 'service/get.event2.data.json',
    }
  );
  cy.intercept('api/controlPlane/v1/project/dynatrace/stage/quality-gate/service/items', {
    fixture: 'get.service.items.mock.json',
  });
}

export function interceptEvaluationBoardWithoutDeployment(): void {
  interceptEvaluationBoard();
  cy.intercept('api/mongodb-datastore/event?keptnContext=*&type=sh.keptn.event.evaluation.triggered&pageSize=1', {
    fixture: 'service/get.evaluation.triggered.mock.json',
  });
}

export function interceptD3(): void {
  cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoEnableD3Heatmap.mock.json' });
}

export function interceptHeatmapComponent(): void {
  interceptD3();
  cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
  cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: false });
  cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' });
  cy.intercept('GET', '/api/project/sockshop/serviceStates', {
    statusCode: 200,
    fixture: 'get.sockshop.service.states.mock.json',
  }).as('serviceStates');
  cy.intercept('GET', '/api/project/sockshop/deployment/da740469-9920-4e0c-b304-0fd4b18d17c2', {
    statusCode: 200,
    fixture: 'get.sockshop.service.carts.deployment.mock.json',
  }).as('ServiceDeployment');
  cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
    statusCode: 200,
    fixture: 'get.sockshop.service.carts.evaluations.heatmap.mock.json',
  }).as('heatmapEvaluations');
}

export function interceptHeatmapComponentWithSLO(slo?: string): void {
  interceptHeatmapComponent();
  cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
    statusCode: 200,
    body: EvaluationFinishedMock(slo),
  }).as('heatmapEvaluations');
}

export function interceptServiceSettings(): void {
  interceptProjectBoard();
  cy.intercept('/api/project/sockshop', { fixture: 'project.mock' }).as('projectPlain');
  cy.intercept('/api/project/sockshop/service/carts/files', {
    body: [],
  });
}

export function interceptHeatmapWithKeySLI(): void {
  interceptHeatmapComponent();
  cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
    statusCode: 200,
    fixture: 'get.sockshop.service.carts.evaluations.keysli.mock.json',
  }).as('heatmapEvaluations');
}

export function interceptHeatmapComponentWithScores(score1: number, score2: number): void {
  interceptHeatmapComponent();
  cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
    statusCode: 200,
    body: EvaluationFinishedScoredMock(score1, score2),
  }).as('heatmapEvaluations');
}
