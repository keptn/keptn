import { init } from '../app';
import { Express } from 'express';
// import { ShipyardResponse } from '../fixtures/shipyard-response';
// eslint-disable-next-line import/no-extraneous-dependencies
let app: Express;

describe('Test the root path', () => {
  beforeAll(async () => {
    app = await init();
  });
  test('It should response the GET method', async () => {
    // axios.get.mockResolvedValue(async () => {
    //   return {};
    // });
    // const response = await request(app).get('/project/sockshop/tasks');
    // expect(response.statusCode).toBe(200);
  });
});
