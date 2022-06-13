import MockAdapter from 'axios-mock-adapter';
import { ApiService } from '../services/api-service';
import { setupServer } from '../.jest/setupServer';
import { AxiosResponse } from 'axios';
import { EnvType } from '../utils/configuration';

describe('Test setting header authorization', () => {
  let apiService: ApiService;
  let axiosMock: MockAdapter;
  const myAccessToken = 'myAccessToken';
  const authorizationHeader = `Bearer ${myAccessToken}`;

  beforeAll(async () => {
    await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  beforeEach(async () => {
    apiService = new ApiService('./', undefined, EnvType.TEST);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should set authorization header if accessToken is provided', async () => {
    await testHeaderResponse(myAccessToken, authorizationHeader);
  });

  it('should not set authorization header if accessToken is not provided', async () => {
    await testHeaderResponse(undefined, undefined);
  });

  async function testHeaderResponse(
    accessToken: string | undefined,
    responseAuthorization: string | undefined
  ): Promise<void> {
    const ignoreMethods = ['getAuthHeaders', 'constructor'];
    axiosMock.onAny().reply(200);

    for (const key of Reflect.ownKeys(Object.getPrototypeOf(apiService))) {
      if (!ignoreMethods.some((method: string) => method === key)) {
        const method = (
          apiService[key as keyof ApiService] as unknown as (
            accessToken: string | undefined,
            ...params: string[]
          ) => Promise<AxiosResponse>
        ).bind(apiService);

        const response = await method(accessToken, '', '', '', '', '', '');
        expect(response.config.headers?.Authorization).toBe(responseAuthorization);
      }
    }
  }
});
