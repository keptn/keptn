import { CallbackParamsType, TokenSet } from 'openid-client';
import request from 'supertest';
import { Express, Request } from 'express';
// eslint-disable-next-line import/no-extraneous-dependencies
import { init } from '../app';

// import { jest } from '@jest/globals';

interface OAuthParameters {
  OAUTH_CLIENT_ID: string | undefined;
  OAUTH_BASE_URL: string | undefined;
  OAUTH_DISCOVERY: string | undefined;
}

const authorizationUrl = `http://localhost/authorization`;
const endSessionEndpoint = 'http://localhost/end_session';
const idToken =
  'myHeader.' +
  Buffer.from(
    JSON.stringify({
      name: 'myName',
    })
  ).toString('base64');

describe('Test OAuth env variables', () => {
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

  it('should not register OAuth endpoints if OAuth is not enabled', async () => {
    delete process.env.OAUTH_ENABLED;
    process.env.OAUTH_CLIENT_ID = 'myClientID';
    process.env.OAUTH_BASE_URL = 'http://localhost';
    process.env.OAUTH_DISCOVERY = 'http://localhost/.well-known/openid-configuration';
    const app = await init();
    for (const endpoint of ['/login', '/oauth/redirect', '/logoutsession']) {
      const response = await request(app).get(endpoint);
      expect(response.status).toBe(500);
    }
    for (const endpoint of ['/logout']) {
      const response = await request(app).post(endpoint);
      expect(response.status).toBe(500);
    }
  });
});

describe('Test OAuth', () => {
  let app: Express;
  beforeAll(async () => {
    mockOpenId(true);
    app = await setupOAuth();
  });

  it('should redirect to authorizationUrl', async () => {
    const response = await request(app).get('/login');
    const state = response.headers.location?.split('state=').pop();
    expect(response.redirect).toBe(true);
    expect(response.headers.location).toEqual(`${authorizationUrl}?state=${state}`);
  });

  it('should reject token gain if state is not provided', async () => {
    await request(app).get('/login');
    const response = await request(app).get(`/oauth/redirect?code=someCode`);
    expect(response.headers.location).toEqual('/');
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
  });

  it('should reject token gain if code is not provided', async () => {
    await request(app).get('/login');
    const response = await request(app).get(`/oauth/redirect?state=someState`);
    expect(response.headers.location).toEqual('/');
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
  });

  it('should reject token gain if state is invalid', async () => {
    await request(app).get('/login');
    const response = await request(app).get(`/oauth/redirect?state=invalidState?code=someCode`);
    expect(response.headers.location).toEqual('/');
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
  });

  it('should authenticate successfully', async () => {
    const { response } = await login(app);
    expect(response.headers['set-cookie']?.[0]?.split('=')[0]).toBe('KTSESSION');
  });

  it('should not be successful if state already used', async () => {
    const { state } = await login(app);
    const response = await request(app).get(`/oauth/redirect?state=${state}&code=someOtherCode`);
    expect(response.headers['set-cookie']?.length ?? 0).toBe(0);
  });

  it('should logout and return end session data', async () => {
    const { response } = await login(app);
    const logoutResponse = await request(app).post(`/logout`).set('Cookie', response.headers['set-cookie']);
    const { state, ...data } = logoutResponse.body;
    expect(state).not.toBeUndefined();
    expect(data).toEqual({
      id_token_hint: idToken,
      post_logout_redirect_uri: 'http://localhost/oauth/redirect',
      end_session_endpoint: endSessionEndpoint,
    });
  });

  it('should return nothing on logout if not authenticated', async () => {
    await login(app);
    const response = await request(app).post(`/logout`);
    expect(response.body).toBe('');
  });

  it('should not be able to fetch data if not authenticated', async () => {
    const response = await request(app).get(`/api/bridgeInfo`);
    expect(response.status).toBe(401);
  });

  it('should be able to fetch data if authenticated', async () => {
    const { response } = await login(app);
    const dataResponse = await request(app).get('/api/bridgeInfo').set('Cookie', response.headers['set-cookie']);
    expect(dataResponse.status).not.toBe(401);
  });
});

describe('Test expired token', () => {
  it('should fail refresh of token and remove session', async () => {
    mockOpenId(true, true, true);
    const app = await setupOAuth();
    const { response } = await login(app);
    const dataResponse = await request(app).get('/api/bridgeInfo').set('Cookie', response.headers['set-cookie']);
    expect(dataResponse.status).toBe(302);
    expect(dataResponse.redirect).toBe(true);
    expect(dataResponse.headers['set-cookie']?.length ?? 0).toBe(0);
  });

  it('should refresh token if expired', async () => {
    mockOpenId(true, true);
    const app = await setupOAuth();
    const { response } = await login(app);
    const dataResponse = await request(app).get('/api/bridgeInfo').set('Cookie', response.headers['set-cookie']);
    expect(dataResponse.status).not.toBe(401);
  });
});

describe('Test OAuth logout without end session endpoint', () => {
  let app: Express;

  beforeAll(async () => {
    mockOpenId(false);
    app = await setupOAuth();
  });

  it('should logout and not return nothing', async () => {
    const { response } = await login(app);
    const logoutResponse = await request(app).post(`/logout`).set('Cookie', response.headers['set-cookie']);
    expect(logoutResponse.body).toBe('');
  });
});

function mockOpenId(includeEndSessionEndpoint: boolean, expiredToken = false, failRefresh = false): void {
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
      return new TokenSet({
        access_token: 'myAccessToken',
        token_type: 'Bearer',
        id_token: idToken,
        refresh_token: 'myRefreshToken',
        scope: 'openid',
        expires_at: new Date().getTime() / 1000 + (expiredToken ? -1 : 10 * 60 * 1000),
      });
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

async function setupOAuth(): Promise<Express> {
  process.env.OAUTH_ENABLED = 'true';
  process.env.OAUTH_CLIENT_ID = 'myClientID';
  process.env.OAUTH_BASE_URL = 'http://localhost';
  process.env.OAUTH_DISCOVERY = 'http://localhost/.well-known/openid-configuration';
  return init();
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

async function login(app: Express): Promise<{ state: string; response: request.Response }> {
  const authUrlResponse = await request(app).get('/login');
  const state = authUrlResponse.headers.location?.split('state=').pop();
  const response = await request(app).get(`/oauth/redirect?state=${state}&code=someCode`);
  return { state, response };
}
