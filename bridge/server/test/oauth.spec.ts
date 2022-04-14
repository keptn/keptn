import { CallbackParamsType, TokenSet } from 'openid-client';
import request from 'supertest';
import { Express, Request } from 'express';
// eslint-disable-next-line import/no-extraneous-dependencies
import { Jest } from '@jest/environment';
// eslint-disable-next-line import/no-extraneous-dependencies
import { jest } from '@jest/globals';

let store: CachedStore = {};
const fakeGetOAuthSecrets = jest.fn();
jest.unstable_mockModule('../user/secrets', () => {
  return {
    getOAuthSecrets: fakeGetOAuthSecrets,
    getOAuthMongoExternalConnectionString(): string {
      return '';
    },
  };
});
// has to be imported after secrets mock
// eslint-disable-next-line @typescript-eslint/naming-convention
const { SessionService } = await import('../user/session');

jest.unstable_mockModule('../user/session', () => {
  return {
    SessionService: jest.fn().mockImplementation(() => {
      return Object.assign(new SessionService(), {
        async saveValidationData(state: string, codeVerifier: string, nonce: string): Promise<void> {
          store[state] = {
            _id: state,
            codeVerifier,
            nonce,
          };
        },
        async getAndRemoveValidationData(state: string): Promise<unknown | undefined> {
          const value = store[state];
          delete store[state];
          return value ? { value } : undefined;
        },
      });
    }),
  };
});

// has to be imported after jest mocked
// eslint-disable-next-line @typescript-eslint/naming-convention
const { TestUtils } = await import('../.jest/test.utils');
const { init } = await import('../app');

// import { jest } from '@jest/globals';

interface OAuthParameters {
  OAUTH_CLIENT_ID: string | undefined;
  OAUTH_BASE_URL: string | undefined;
  OAUTH_DISCOVERY: string | undefined;
}

interface CachedStore {
  [state: string]: { _id: string; codeVerifier: string; nonce: string };
}

const authorizationUrl = `http://localhost/authorization`;
const endSessionEndpoint = 'http://localhost/end_session';
let sessionMock: Jest | undefined;
const idToken =
  'myHeader.' +
  Buffer.from(
    JSON.stringify({
      name: 'myName',
    })
  ).toString('base64');

describe('Test OAuth env variables', () => {
  beforeEach(() => {
    mockSecrets();
  });

  it('should exit if insufficient parameters', async () => {
    const parameters: OAuthParameters[] = [
      {
        OAUTH_CLIENT_ID: 'myClientID',
        OAUTH_BASE_URL: 'http://localhost',
        OAUTH_DISCOVERY: undefined,
      },
      {
        OAUTH_CLIENT_ID: 'myClientID',
        OAUTH_BASE_URL: undefined,
        OAUTH_DISCOVERY: 'http://localhost/.well-known/openid-configuration',
      },
      {
        OAUTH_CLIENT_ID: undefined,
        OAUTH_BASE_URL: 'http://localhost',
        OAUTH_DISCOVERY: 'http://localhost/.well-known/openid-configuration',
      },
    ];
    for (const parameter of parameters) {
      process.env.OAUTH_ENABLED = 'true';
      setOrDeleteProperty(process.env, parameter, 'OAUTH_CLIENT_ID');
      setOrDeleteProperty(process.env, parameter, 'OAUTH_BASE_URL');
      setOrDeleteProperty(process.env, parameter, 'OAUTH_DISCOVERY');

      await expect(init()).rejects.toThrowError();
    }
  });

  it('should throw errors if session secret is not provided', async () => {
    process.env.OAUTH_ENABLED = 'true';
    process.env.OAUTH_CLIENT_ID = 'myClientID';
    process.env.OAUTH_BASE_URL = 'http://localhost';
    process.env.OAUTH_DISCOVERY = 'http://localhost/.well-known/openid-configuration';
    fakeGetOAuthSecrets.mockImplementation(() => {
      return {
        sessionSecret: '',
        databaseEncryptSecret: 'database_secret_'.repeat(2),
      };
    });
    await expect(init()).rejects.toThrowError();
  });

  it('should throw errors if database encrypt secret is not provided', async () => {
    process.env.OAUTH_ENABLED = 'true';
    process.env.OAUTH_CLIENT_ID = 'myClientID';
    process.env.OAUTH_BASE_URL = 'http://localhost';
    process.env.OAUTH_DISCOVERY = 'http://localhost/.well-known/openid-configuration';
    fakeGetOAuthSecrets.mockImplementation(() => {
      return {
        sessionSecret: 'abcd',
        databaseEncryptSecret: '',
      };
    });
    await expect(init()).rejects.toThrowError();
  });

  it('should throw errors if database encrypt secret length is invalid', async () => {
    process.env.OAUTH_ENABLED = 'true';
    process.env.OAUTH_CLIENT_ID = 'myClientID';
    process.env.OAUTH_BASE_URL = 'http://localhost';
    process.env.OAUTH_DISCOVERY = 'http://localhost/.well-known/openid-configuration';
    fakeGetOAuthSecrets.mockImplementation(() => {
      return {
        sessionSecret: 'abcd',
        databaseEncryptSecret: 'mySecret',
      };
    });
    await expect(init()).rejects.toThrowError();
  });

  it('should not register OAuth endpoints if OAuth is not enabled', async () => {
    delete process.env.OAUTH_ENABLED;
    process.env.OAUTH_CLIENT_ID = 'myClientID';
    process.env.OAUTH_BASE_URL = 'http://localhost';
    process.env.OAUTH_DISCOVERY = 'http://localhost/.well-known/openid-configuration';
    const app = await init();
    // if not found, index.html is sent
    for (const endpoint of ['/oauth/login', '/oauth/redirect']) {
      const response = await request(app).get(endpoint);
      expect(response.text).not.toBeUndefined();
    }
    for (const endpoint of ['/oauth/logout']) {
      const response = await request(app).post(endpoint);
      expect(response.text).not.toBeUndefined();
    }
  });

  afterAll(() => {
    sessionMock?.resetAllMocks();
  });
});

describe('Test OAuth', () => {
  let app: Express;
  beforeAll(async () => {
    mockOpenId(true);
    app = await setupOAuth();
  });

  beforeEach(() => {
    store = {};
  });

  it('should redirect to authorizationUrl', async () => {
    const response = await request(app).get('/oauth/login/');
    const state = response.headers.location?.split('state=').pop();
    expect(response.redirect).toBe(true);
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
    expect(response.headers.location).toEqual(`${authorizationUrl}?state=${state}`);
  });

  it('should reject token gain if state is not provided', async () => {
    await request(app).get('/oauth/login/');
    const response = await request(app).get(`/oauth/redirect?code=someCode`);
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
    expect(response.headers.location).toEqual('/error');
    expect(response.text).not.toBeUndefined(); // error view is rendered
  });

  it('should reject token gain if code is not provided', async () => {
    await request(app).get('/oauth/login/');
    const response = await request(app).get(`/oauth/redirect?state=someState`);
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
    expect(response.headers.location).toEqual('/error');
    expect(response.text).not.toBeUndefined();
  });

  it('should reject token gain if state is invalid', async () => {
    await request(app).get('/oauth/login/');
    const response = await request(app).get(`/oauth/redirect?state=invalidState?code=someCode`);
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
    expect(response.headers.location).toEqual('/error');
    expect(response.text).not.toBeUndefined();
  });

  it('should authenticate successfully', async () => {
    const { response, cookies } = await login(app);
    expect(cookies?.[0]?.split('=')[0]).toBe('KTSESSION');
    expect(response.headers.location).toBe('/');
    expect(response.status).toBe(302);
  });

  it('should not be successful if state already used', async () => {
    const { state } = await login(app);
    const response = await request(app).get(`/oauth/redirect?state=${state}&code=someOtherCode`);
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
    expect(response.headers.location).toEqual('/error');
    expect(response.text).not.toBeUndefined();
  });

  it('should logout and return end session data', async () => {
    const { cookies } = await login(app);
    const logoutResponse = await request(app).post(`/oauth/logout`).set('Cookie', cookies);
    const { state, ...data } = logoutResponse.body;
    expect(state).not.toBeUndefined();
    expect(data).toEqual({
      id_token_hint: idToken,
      post_logout_redirect_uri: 'http://localhost/logoutsession',
      end_session_endpoint: endSessionEndpoint,
    });
  });

  it('should return nothing on logout if not authenticated', async () => {
    const response = await request(app).post(`/oauth/logout`);
    expect(response.body).toBe('');
  });

  it('should not be able to fetch data if not authenticated', async () => {
    const response = await request(app).get(`/api/bridgeInfo`);
    expect(response.status).toBe(401);
  });

  it('should be able to fetch data if authenticated', async () => {
    const { cookies } = await login(app);
    const dataResponse = await request(app).get('/api/bridgeInfo').set('Cookie', cookies);
    expect(dataResponse.status).not.toBe(401);
  });

  afterAll(() => {
    sessionMock?.resetAllMocks();
  });
});

describe('Test OAuth with error on callback', () => {
  beforeEach(() => {
    store = {};
  });

  it('should redirect to insufficient permission page if callback fails', async () => {
    mockOpenId(true, false, false, 403);
    const app = await setupOAuth();
    const { response } = await login(app);
    expect(response.headers.location).toEqual('/error?status=403');
  });

  it('should redirect to internal error page if callback fails', async () => {
    mockOpenId(true, false, false, 500);
    const app = await setupOAuth();
    const { response } = await login(app);
    expect(response.headers.location).toEqual('/error');
  });

  afterAll(() => {
    sessionMock?.resetAllMocks();
  });
});

describe('Test expired token', () => {
  beforeEach(() => {
    store = {};
  });

  it('should fail refresh of token and remove session', async () => {
    mockOpenId(true, true, true);
    const app = await setupOAuth();
    const { cookies } = await login(app);
    const dataResponse = await request(app).get('/api/bridgeInfo').set('Cookie', cookies);

    expect(dataResponse.status).toBe(302);
    expect(dataResponse.redirect).toBe(true);
    expect(dataResponse.headers['set-cookie']?.length ?? 0).toBe(0);
  });

  it('should refresh token if expired', async () => {
    mockOpenId(true, true);
    const app = await setupOAuth();
    const { cookies } = await login(app);
    const dataResponse = await request(app).get('/api/bridgeInfo').set('Cookie', cookies);

    expect(dataResponse.status).not.toBe(401);
  });

  afterAll(() => {
    sessionMock?.resetAllMocks();
  });
});

describe('Test OAuth logout without end session endpoint', () => {
  let app: Express;

  beforeAll(async () => {
    mockOpenId(false);
    app = await setupOAuth();
  });

  beforeEach(() => {
    store = {};
  });

  it('should logout and not return nothing', async () => {
    const { cookies } = await login(app);
    const logoutResponse = await request(app).post(`/oauth/logout`).set('Cookie', cookies);
    expect(logoutResponse.body).toBe('');
  });

  afterAll(() => {
    sessionMock?.resetAllMocks();
  });
});

function mockOpenId(
  includeEndSessionEndpoint: boolean,
  expiredToken = false,
  failRefresh = false,
  exceptionStatusCodeOnCallback?: number
): void {
  // jest currently does not really support mocking of ESM

  const issuer = {
    metadata: {
      issuer: 'http://localhost',
      authorization_endpoint: authorizationUrl,
      token_endpoint: 'http://localhost/token',
      ...(includeEndSessionEndpoint && { end_session_endpoint: endSessionEndpoint }),
    },
  };

  class MockBaseClient {
    issuer = issuer;
    callbackParams(req: Request): CallbackParamsType {
      return {
        state: req.query.state?.toString(),
        code: req.query.code?.toString(),
      };
    }
    async callback(): Promise<TokenSet> {
      if (exceptionStatusCodeOnCallback !== undefined) {
        const error = new Error() as Error & { response: { statusCode: number } };

        error.response = {
          statusCode: exceptionStatusCodeOnCallback,
        };
        throw error;
      } else {
        return new TokenSet({
          access_token: 'myAccessToken',
          token_type: 'Bearer',
          id_token: idToken,
          refresh_token: 'myRefreshToken',
          scope: 'openid',
          expires_at: new Date().getTime() / 1000 + (expiredToken ? -1 : 10 * 60 * 1000),
        });
      }
    }

    async refresh(tokenSet: TokenSet): Promise<TokenSet> {
      if (failRefresh) {
        throw new Error('Refresh failed');
      }
      tokenSet.expires_at = new Date().getTime() / 1000 + 10 * 60 * 1000;
      return tokenSet;
    }

    authorizationUrl({ state }: { state: string }): string {
      return `${authorizationUrl}?state=${state}`;
    }
  }

  global.issuer = {
    discover(): Promise<unknown> {
      return new Promise((resolve) => {
        resolve({ ...issuer, Client: MockBaseClient });
      });
    },
  };
}

function setOrDeleteProperty(
  env: Record<string, string | undefined>,
  parameter: OAuthParameters,
  key: keyof OAuthParameters
): void {
  if (parameter[key]) {
    process.env[key] = parameter[key];
  } else {
    delete process.env[key];
  }
}

async function login(app: Express): Promise<{ state: string; response: request.Response; cookies: string[] }> {
  const authUrlResponse = await request(app).get('/oauth/login/');
  const state = authUrlResponse.headers.location?.split('state=').pop();
  const response = await request(app).get(`/oauth/redirect?state=${state}&code=someCode`);
  return { state, response, cookies: response.headers['set-cookie'] };
}

async function setupOAuth(): Promise<Express> {
  mockSecrets();
  await mockSavingValidationData();
  return TestUtils.setupOAuthTest();
}

async function mockSavingValidationData(): Promise<void> {
  process.env.OAUTH_SESSION_SECRET = 'mySessionSecret';
  process.env.OAUTH_DATABASE_ENCRYPT_SECRET = 'database_secret_'.repeat(2); // length of 32
}

function mockSecrets(): void {
  fakeGetOAuthSecrets.mockImplementation(() => {
    return {
      sessionSecret: 'abc',
      databaseEncryptSecret: 'database_secret_'.repeat(2),
    };
  });
}
