import MockAdapter from 'axios-mock-adapter';
import request from 'supertest';
import { ProjectResponseIntersect } from '../../shared/fixtures/project-response.mock';
import { EventTypes } from '../../shared/interfaces/event-types';
import {
  IntersectDeploymentFinishedResponse,
  IntersectDeploymentStartedResponse,
  IntersectDeploymentTriggeredResponse,
} from '../../shared/fixtures/traces-response.mock';
import { IntersectEventResponse } from '../fixtures/intersect-event-response.mock';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

let axiosMock: MockAdapter;

describe('Test /intersectEvents', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve intersection of stage and service and one event', async () => {
    const projectName = 'sockshop';
    const event1 = {
      data: {
        project: 'sockshop',
        service: 'carts',
        stage: 'dev',
      },
      id: 'myId',
      type: EventTypes.DEPLOYMENT_FINISHED,
      keptnContext: 'myContext',
    };
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}`).reply(200, ProjectResponseIntersect);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND id:3`,
          limit: '100',
        },
      })
      .reply(200, {
        events: [event1],
      });
    const response = await request(app)
      .post(`/api/intersectEvents`)
      .send({
        projectName: 'sockshop',
        stages: ['dev'],
        services: ['carts'],
        event: 'sh.keptn.event.deployment',
        eventSuffix: 'finished',
      });
    expect(response.status).toBe(200);
    expect(response.body).toEqual(event1);
  });

  it('should retrieve intersection of multiple stage, services and events', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}`).reply(200, ProjectResponseIntersect);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_TRIGGERED}`, {
        params: {
          filter: `data.project:${projectName} AND id:1,4,7,10`,
          limit: '100',
        },
      })
      .reply(200, IntersectDeploymentTriggeredResponse);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_STARTED}`, {
        params: {
          filter: `data.project:${projectName} AND id:2,5,8,11`,
          limit: '100',
        },
      })
      .reply(200, IntersectDeploymentStartedResponse);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND id:3,6,9`,
          limit: '100',
        },
      })
      .reply(200, IntersectDeploymentFinishedResponse);
    const response = await request(app).post(`/api/intersectEvents`).send({
      projectName: 'sockshop',
      stages: [],
      services: [],
      event: 'sh.keptn.event.deployment',
      eventSuffix: '>',
    });
    expect(response.status).toBe(200);
    expect(response.body).toEqual(IntersectEventResponse);
  });

  it('should correctly intersect objects and arrays', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}`).reply(200, ProjectResponseIntersect);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_TRIGGERED}`, {
        params: {
          filter: `data.project:${projectName} AND id:1`,
          limit: '100',
        },
      })
      .reply(200, {
        events: [
          {
            data: {
              myCustomData: [
                [
                  {
                    element1: 'element1',
                    elementZ: 'elementZ',
                    elementType: 1,
                    elementProps: [],
                  },
                ],
                {
                  element2: 'element2',
                },
                'myString',
                null,
              ],
            },
            prop: null,
            id: 'id',
            customProperty: 'customProperty',
          },
          {
            data: {
              myCustomData: [
                [
                  {
                    element1: 'element1',
                    elementX: 'elementX',
                    elementType: '1',
                    elementProps: {},
                  },
                ],
                {
                  element2: 'element2',
                  elementY: 'elementY',
                },
                'myString2',
                'myString3',
                {
                  additionalElement: 'additionalElement',
                },
              ],
            },
            id: 'id',
          },
        ],
      });
    const response = await request(app)
      .post(`/api/intersectEvents`)
      .send({
        projectName: 'sockshop',
        stages: ['dev'],
        services: ['carts'],
        event: 'sh.keptn.event.deployment',
        eventSuffix: 'triggered',
      });
    expect(response.status).toBe(200);
    expect(response.body).toEqual({
      data: {
        myCustomData: [
          [
            {
              element1: 'element1',
              elementType: 1,
            },
          ],
          {
            element2: 'element2',
          },
          'myString',
          null,
        ],
      },
      id: 'id',
    });
  });

  it('should return empty object if there are no events for intersection', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}`).reply(200, ProjectResponseIntersect);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.DEPLOYMENT_TRIGGERED}`, {
        params: {
          filter: `data.project:${projectName} AND id:1`,
          limit: '100',
        },
      })
      .reply(200, {
        events: [],
      });
    const response = await request(app)
      .post(`/api/intersectEvents`)
      .send({
        projectName: 'sockshop',
        stages: ['dev'],
        services: ['carts'],
        event: 'sh.keptn.event.deployment',
        eventSuffix: 'triggered',
      });
    expect(response.status).toBe(200);
    expect(response.body).toEqual({});
  });

  it('should send error if projectName is missing', async () => {
    const response = await request(app).post(`/api/intersectEvents`).send({
      event: 'sh.keptn.event.deployment',
      eventSuffix: 'finished',
    });
    expect(response.status).toBe(400);
  });

  it('should send error if event is missing', async () => {
    const response = await request(app).post(`/api/intersectEvents`).send({
      projectName: 'sockshop',
      eventSuffix: 'finished',
    });
    expect(response.status).toBe(400);
  });

  it('should send error if eventSuffix is missing', async () => {
    const response = await request(app).post(`/api/intersectEvents`).send({
      projectName: 'sockshop',
      event: 'sh.keptn.event.deployment',
    });
    expect(response.status).toBe(400);
  });
});
