import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { StagesResponse } from '../fixtures/stages.mock';
import {
  ProjectDetailsResponseEvaluationFallback,
  ProjectDetailsResponseURLFallback,
  ProjectResponse,
  ProjectResponseEvaluationFallback,
  ProjectResponseURLFallback,
} from '../../shared/fixtures/project-response.mock';
import { EventTypes } from '../../shared/interfaces/event-types';
import { SequenceState } from '../../shared/models/sequence';
import {
  OpenRemediationsResponse,
  RemediationTracesResponse,
} from '../../shared/fixtures/open-remediations-response.mock';
import {
  ApprovalEvaluationResponse,
  DefaultDeploymentData,
  DefaultDeploymentFinishedTrace,
  DefaultEvaluationFinishedTrace,
  LatestFinishedDeployments,
  LatestFinishedEvaluations,
  OpenApprovalsResponse,
} from '../../shared/fixtures/traces-response.mock';
import {
  SequenceResponseEvaluationFallback,
  SequenceResponseURLFallback,
  SequencesResponses,
} from '../fixtures/sequence-response.mock';
import { KeptnService } from '../../shared/models/keptn-service';
import { ProjectDetailsResponse } from '../fixtures/project-details-response.mock';
import { ResultTypes } from '../../shared/models/result-types';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

let axiosMock: MockAdapter;

describe('Test project resources', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve service names', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}/stage`).reply(200, StagesResponse);
    const response = await request(app).get(`/api/project/${projectName}/services`);
    expect(response.body).toEqual(['carts', 'carts-db']);
    expect(response.statusCode).toBe(200);
  });

  it('should return an error', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}/stage`).reply(502);
    const response = await request(app).get(`/api/project/${projectName}/services`);
    expect(response.statusCode).toBe(502);
  });

  it('should fetch and aggregate project details', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}`).reply(200, ProjectResponse);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_FINISHED}`)
      .reply(200, LatestFinishedDeployments);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`)
      .reply(200, LatestFinishedEvaluations);
    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
        params: {
          pageSize: '100',
          name: 'remediation',
          state: SequenceState.STARTED,
        },
      })
      .reply(200, OpenRemediationsResponse);
    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/event/triggered/${EventTypes.APPROVAL_TRIGGERED}`, {
        params: {
          project: projectName,
        },
      })
      .reply(200, OpenApprovalsResponse);

    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event`, {
        params: {
          project: 'sockshop',
          service: 'carts',
          stage: 'production',
          keptnContext: '35383737-3630-4639-b037-353138323631',
          pageSize: '50',
        },
      })
      .reply(200, RemediationTracesResponse);

    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event`, {
        params: {
          project: 'sockshop',
          type: EventTypes.EVALUATION_FINISHED,
          pageSize: '1',
          keptnContext: OpenApprovalsResponse.events[0].shkeptncontext,
          source: KeptnService.LIGHTHOUSE_SERVICE,
        },
      })
      .reply(200, ApprovalEvaluationResponse);

    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`).reply((config) => {
      const context = config.params.keptnContext;
      const sequence = SequencesResponses[context];
      expect(sequence).not.toBeUndefined();
      return [200, sequence];
    });

    const response = await request(app).get(`/api/project/${projectName}?approval=true&remediation=true`);
    expect(response.body).toEqual(ProjectDetailsResponse);
  });

  it('should correctly fallback to right deployment URL', async () => {
    const projectName = 'sockshop';
    const data = {
      message: 'Failed to deploy',
      project: 'sockshop',
      result: ResultTypes.FAILED,
      service: 'carts',
      stage: 'dev',
      status: 'failed',
    };
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}`).reply(200, ProjectResponseURLFallback);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND id:eventId`,
          excludeInvalidated: 'true',
        },
      })
      .reply(200, {
        events: [
          {
            ...DefaultDeploymentFinishedTrace,
            data,
          },
        ],
      });

    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND data.service:carts AND data.stage:dev AND data.result:pass`,
          excludeInvalidated: 'true',
          limit: '1',
        },
      })
      .reply(200, {
        events: [
          {
            ...DefaultDeploymentFinishedTrace,
            data: {
              ...data,
              deployment: DefaultDeploymentData,
              result: ResultTypes.PASSED,
            },
          },
        ],
      });

    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`)
      .reply(200, SequenceResponseURLFallback);

    const response = await request(app).get(`/api/project/${projectName}`);
    expect(response.body).toEqual(ProjectDetailsResponseURLFallback);
  });

  it('should correctly fallback to right evaluation trace', async () => {
    const projectName = 'sockshop';

    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}`)
      .reply(200, ProjectResponseEvaluationFallback);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND id:webhookId`,
          excludeInvalidated: 'true',
        },
      })
      .reply(200, {
        events: [
          {
            ...DefaultEvaluationFinishedTrace,
            source: 'webhook-service',
          },
        ],
      });

    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND data.service:carts AND data.stage:dev AND data.source:${KeptnService.LIGHTHOUSE_SERVICE}`,
          excludeInvalidated: 'true',
          limit: '1',
        },
      })
      .reply(200, {
        events: [
          {
            ...DefaultEvaluationFinishedTrace,
          },
        ],
      });

    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`)
      .reply(200, SequenceResponseEvaluationFallback);

    const response = await request(app).get(`/api/project/${projectName}`);
    expect(response.body).toEqual(ProjectDetailsResponseEvaluationFallback);
  });
});
