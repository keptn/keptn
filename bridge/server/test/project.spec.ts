import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { StagesResponse } from '../fixtures/stages';

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
    axiosMock.onGet(global.baseUrl + `/controlPlane/v1/project/${projectName}/stage`).reply(200, StagesResponse);
    const response = await request(global.app).get(`/api/project/${projectName}/services`);
    expect(response.body).toEqual(['carts', 'carts-db']);
    expect(response.statusCode).toBe(200);
  });
});
