import { BridgeConfiguration, EnvType } from '../interfaces/configuration';
import { EnvVar, getConfiguration } from './configuration';
import { LogDestination } from './logger';

describe('Configuration', () => {
  beforeEach(() => {
    cleanEnv();
    process.env[EnvVar.API_TOKEN] = 'some value to do not run kubectl cmds as part of the tests';
  });
  afterEach(() => {
    cleanEnv();
  });

  function cleanEnv(): void {
    for (const e in EnvVar) {
      delete process.env[e];
    }
  }

  const defaultAPIURL = 'http://localhost';
  const defaultAPIToken = 'abcdefgh';

  function setBasicEnvVar(): void {
    process.env[EnvVar.API_URL] = defaultAPIURL;
    process.env[EnvVar.API_TOKEN] = defaultAPIToken;
    process.env[EnvVar.MONGODB_HOST] = '';
    process.env[EnvVar.MONGODB_USER] = '';
    process.env[EnvVar.MONGODB_PASSWORD] = '';
  }

  it('should use default values', () => {
    setBasicEnvVar();
    const result = getConfiguration();
    const expected: BridgeConfiguration = {
      logging: {
        destination: LogDestination.STDOUT,
        enabledComponents: {},
      },
      api: {
        showToken: true,
        token: defaultAPIToken,
        url: defaultAPIURL,
      },
      auth: {
        authMessage: `keptn auth --endpoint=${defaultAPIURL} --api-token=${defaultAPIToken}`,
        basicPassword: undefined,
        basicUsername: undefined,
        cleanBucketIntervalMs: 60 * 60 * 1000, //1h
        requestTimeLimitMs: 60 * 60 * 1000, //1h
        nRequestWithinTime: 10,
      },
      oauth: {
        allowedLogoutURL: '',
        baseURL: '',
        clientID: '',
        clientSecret: undefined,
        discoveryURL: '',
        enabled: false,
        nameProperty: undefined,
        scope: '',
        session: {
          secureCookie: false,
          timeoutMin: 60,
          trustProxyHops: 1,
          validationTimeoutMin: 60,
        },
        tokenAlgorithm: 'RS256',
      },
      urls: {
        CLI: 'https://github.com/keptn/keptn/releases',
        integrationPage: 'https://get.keptn.sh/integrations.html',
        lookAndFeel: undefined,
      },
      features: {
        automaticProvisioningMessage: '',
        configDir: 'config',
        installationType: 'QUALITY_GATES,CONTINUOUS_OPERATIONS,CONTINUOUS_DELIVERY',
        pageSize: {
          project: 50,
          service: 50,
        },
        prefixPath: '/',
        versionCheck: true,
      },
      mongo: {
        db: 'openid',
        host: '',
        password: '',
        user: '',
      },
      mode: EnvType.DEV,
      version: 'develop',
    };

    // handlign special cases separately and then patch them to known values
    expect(result.features.configDir).toMatch(/config$/);
    result.features.configDir = 'config';

    expect(result).toStrictEqual(expected);
  });

  it('should set values using options object', () => {
    setBasicEnvVar();
    const oauthBaseUrl = 'mybaseurl';
    const oauthClientID = 'myclientid';
    const oauthDiscovery = 'mydiscovery';
    const apiUrl = 'myapiurl';
    const apiToken = 'mytoken';
    const version = '0.0.0';
    const mongoHost = 'localhost';
    let result = getConfiguration({
      logging: {
        enabledComponents: 'a=true,b=false,c=true',
      },
      oauth: {
        baseURL: oauthBaseUrl,
        clientID: oauthClientID,
        discoveryURL: oauthDiscovery,
        enabled: true,
      },
      api: {
        showToken: false,
        token: apiToken,
        url: apiUrl,
      },
      mongo: {
        user: '',
        password: '',
        host: mongoHost,
      },
      mode: 'test',
      version: version,
    });

    const expected: BridgeConfiguration = {
      logging: {
        destination: LogDestination.STDOUT,
        enabledComponents: {
          a: true,
          b: false,
          c: true,
        },
      },
      api: {
        showToken: false,
        token: apiToken,
        url: apiUrl,
      },
      auth: {
        authMessage: `keptn auth --endpoint=${apiUrl} --api-token=${apiToken}`,
        basicPassword: undefined,
        basicUsername: undefined,
        cleanBucketIntervalMs: 60 * 60 * 1000, //1h
        requestTimeLimitMs: 60 * 60 * 1000, //1h
        nRequestWithinTime: 10,
      },
      oauth: {
        allowedLogoutURL: '',
        baseURL: oauthBaseUrl,
        clientID: oauthClientID,
        clientSecret: undefined,
        discoveryURL: oauthDiscovery,
        enabled: true,
        nameProperty: undefined,
        scope: '',
        session: {
          secureCookie: false,
          timeoutMin: 60,
          trustProxyHops: 1,
          validationTimeoutMin: 60,
        },
        tokenAlgorithm: 'RS256',
      },
      urls: {
        CLI: 'https://github.com/keptn/keptn/releases',
        integrationPage: 'https://get.keptn.sh/integrations.html',
        lookAndFeel: undefined,
      },
      features: {
        automaticProvisioningMessage: '',
        configDir: 'config',
        installationType: 'QUALITY_GATES,CONTINUOUS_OPERATIONS,CONTINUOUS_DELIVERY',
        pageSize: {
          project: 50,
          service: 50,
        },
        prefixPath: '/',
        versionCheck: true,
      },
      mongo: {
        db: 'openid',
        host: mongoHost,
        password: '',
        user: '',
      },
      mode: EnvType.TEST,
      version: version,
    };

    // handlign special cases separately and then patch them to known values
    expect(result.features.configDir).toMatch(/config$/);
    result.features.configDir = 'config';

    expect(result).toStrictEqual(expected);

    // check that values can change
    result = getConfiguration({
      logging: {
        enabledComponents: 'a=false',
        destination: LogDestination.FILE,
      },
    });
    expect(result.logging.destination).toBe(LogDestination.FILE);
    expect(result.logging.enabledComponents).toStrictEqual({
      a: false,
    });
  });

  it('should fail for missing API values', () => {
    expect(getConfiguration).toThrow('API_URL is not provided');
  });

  it('should fail for missing OAuth values', () => {
    process.env[EnvVar.MONGODB_HOST] = 'mongo://';
    process.env[EnvVar.MONGODB_PASSWORD] = 'pwd';
    process.env[EnvVar.MONGODB_USER] = 'usr';
    expect(() => {
      getConfiguration({
        api: { url: 'somevalue' },
        oauth: { enabled: true },
      });
    }).toThrow(/OAUTH_.*/);
    process.env[EnvVar.OAUTH_ENABLED] = 'true';
    const t: () => void = () => {
      getConfiguration({
        api: { url: 'somevalue' },
      });
    };
    expect(t).toThrow(/OAUTH_.*/);
    process.env[EnvVar.OAUTH_DISCOVERY] = 'http://keptn';
    expect(t).toThrow(/OAUTH_.*/);
    process.env[EnvVar.OAUTH_CLIENT_ID] = 'abcdefg';
    expect(t).toThrow(/OAUTH_.*/);
    process.env[EnvVar.OAUTH_BASE_URL] = 'http://keptn';
    expect(t).not.toThrow();
  });

  it('should fail for missing Mongo values', () => {
    process.env[EnvVar.API_URL] = 'http://localhost';
    process.env[EnvVar.OAUTH_ENABLED] = 'true';
    process.env[EnvVar.OAUTH_DISCOVERY] = 'smth';
    process.env[EnvVar.OAUTH_CLIENT_ID] = 'id';
    process.env[EnvVar.OAUTH_BASE_URL] = 'url';
    expect(getConfiguration).toThrow(/Could not construct mongodb connection string.*/);
    process.env[EnvVar.MONGODB_HOST] = 'mongo://';
    expect(getConfiguration).toThrow(/Could not construct mongodb connection string.*/);
    process.env[EnvVar.MONGODB_PASSWORD] = 'pwd';
    expect(getConfiguration).toThrow(/Could not construct mongodb connection string.*/);
    process.env[EnvVar.MONGODB_USER] = 'usr';
    expect(getConfiguration).not.toThrow();
  });

  it('should set values using env var', () => {
    setBasicEnvVar();
    process.env.LOGGING_COMPONENTS = 'a=true,b=false,c=true';
    const result = getConfiguration();
    expect(result.logging.destination).toBe(LogDestination.STDOUT);
    expect(result.logging.enabledComponents).toStrictEqual({
      a: true,
      b: false,
      c: true,
    });
  });

  it('option object should win over env var', () => {
    setBasicEnvVar();
    process.env.LOGGING_COMPONENTS = 'a=false,b=true,c=false';
    const result = getConfiguration({
      logging: {
        enabledComponents: 'a=true,b=false,c=true',
      },
    });
    expect(result.logging.destination).toBe(LogDestination.STDOUT);
    expect(result.logging.enabledComponents).toStrictEqual({
      a: true,
      b: false,
      c: true,
    });
  });

  it('should correctly eval booleans', () => {
    setBasicEnvVar();
    const result = getConfiguration({
      logging: {
        enabledComponents: 'a=tRue,b=FaLsE,c=0,d=1,e=enabled',
      },
    });
    expect(result.logging.enabledComponents).toStrictEqual({
      a: true,
      b: false,
      c: false,
      d: true,
      e: true,
    });
  });
});
