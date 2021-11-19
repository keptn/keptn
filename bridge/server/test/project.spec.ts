import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { StagesResponse } from '../fixtures/stages.mock';
import { ProjectResponse } from '../fixtures/project-response.mock';
import { EventTypes } from '../../shared/interfaces/event-types';
import { SequenceState } from '../../shared/models/sequence';
import {
  OpenRemediationsResponse,
  RemediationTracesResponse,
} from '../../shared/fixtures/open-remediations-response.mock';
import {
  ApprovalEvaluationResponse,
  LatestFinishedDeployments,
  LatestFinishedEvaluations,
  OpenApprovalsResponse,
} from '../fixtures/traces-response.mock';
import { SequencesResponses } from '../fixtures/sequence-response.mock';
import { KeptnService } from '../../shared/models/keptn-service';
import { ProjectDetailsResponse } from '../fixtures/project-details-response.mock';

let axiosMock: MockAdapter;

describe('Test project resources', () => {
  beforeAll(() => {
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve service names', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}/stage`).reply(200, StagesResponse);
    const response = await request(global.app).get(`/api/project/${projectName}/services`);
    expect(response.body).toEqual(['carts', 'carts-db']);
    expect(response.statusCode).toBe(200);
  });

  it('should return an error', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}/stage`).reply(502);
    const response = await request(global.app).get(`/api/project/${projectName}/services`);
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

    const response = await request(global.app).get(`/api/project/${projectName}?approval=true&remediation=true`);
    expect(response.body).toEqual(ProjectDetailsResponse);
  });
});
