import {
  BridgeConfiguration,
  BridgeOption,
  EnvType,
  EnvVar,
  MongoConfig,
  OAuthConfig,
} from '../interfaces/configuration';
import { LogDestination } from './logger';
// eslint-disable-next-line import/no-extraneous-dependencies
import { jest } from '@jest/globals';

const getOAuthSecretsSpy = jest.fn();

jest.unstable_mockModule('../user/secrets', () => {
  return {
    getOAuthSecrets: getOAuthSecretsSpy,
    getOAuthMongoExternalConnectionString: (): string => '',
  };
});

const { envToConfiguration, getConfiguration } = await import('./configuration');

describe('Configuration', () => {
  const defaultAPIURL = 'http://localhost';
  const defaultAPIToken = 'abcdefgh';

  beforeEach(() => {
    getOAuthSecretsSpy.mockReturnValue({
      sessionSecret: 'session_secret',
      databaseEncryptSecret: 'database_secret_'.repeat(2),
      clientSecret: '',
    });
  });

  it('should use default values', () => {
    const basicEnv = getBasicEnvVar();
    const result = getConfiguration(envToConfiguration(basicEnv));
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
        secrets: {
          sessionSecret: '',
          databaseEncryptSecret: '',
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
        externalConnectionString: '',
        host: '',
        password: '',
        user: '',
      },
      mode: EnvType.TEST,
      version: 'develop',
    };

    // handlign special cases separately and then patch them to known values
    expect(result.features.configDir).toMatch(/config$/);
    result.features.configDir = 'config';

    expect(result).toStrictEqual(expected);
  });

  it('should set values using options object', () => {
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
        externalConnectionString: '',
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
        secrets: {
          clientSecret: '',
          sessionSecret: 'session_secret',
          databaseEncryptSecret: 'database_secret_'.repeat(2),
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
        externalConnectionString: '',
        host: mongoHost,
        password: '',
        user: '',
      },
      mode: EnvType.TEST,
      version: version,
    };

    // handling special cases separately and then patch them to known values
    expect(result.features.configDir).toMatch(/config$/);
    result.features.configDir = 'config';

    expect(result).toStrictEqual(expected);

    // check that values can change
    result = getConfiguration({
      ...envToConfiguration(getBasicEnvVar()),
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
    expect(() => getConfiguration({})).toThrow('API_URL is not provided');
  });

  it('should fail for missing OAuth values', () => {
    const config: Omit<BridgeOption, 'oauth'> & { oauth: Partial<OAuthConfig> } = {
      api: { url: 'somevalue' },
      oauth: { enabled: true },
      mode: 'test',
      mongo: {
        host: 'mongo://',
        password: 'pwd',
        user: 'usr',
      },
    };
    expect(() => {
      getConfiguration({
        api: { url: 'somevalue' },
        oauth: { enabled: true },
        mode: 'test',
        mongo: {
          host: 'mongo://',
          password: 'pwd',
          user: 'usr',
        },
      });
    }).toThrow(/OAUTH_.*/);

    const t: () => void = () => {
      getConfiguration(config);
    };
    expect(t).toThrow(/OAUTH_.*/);
    config.oauth.discoveryURL = 'http://keptn';
    expect(t).toThrow(/OAUTH_.*/);
    config.oauth.clientID = 'abcdefg';
    expect(t).toThrow(/OAUTH_.*/);
    config.oauth.baseURL = 'http://keptn';
    expect(t).not.toThrow();
  });

  it('should throw error if OAuth secrets are not provided or invalid', () => {
    const config = envToConfiguration(getDefaultOAuthConfig());
    getOAuthSecretsSpy.mockReturnValue({
      sessionSecret: '',
      databaseEncryptSecret: '',
      clientSecret: '',
    });
    const t = (): BridgeConfiguration => getConfiguration(config);
    expect(t).toThrow(/^session_secret.*/);

    getOAuthSecretsSpy.mockReturnValue({
      sessionSecret: 'asdf',
      databaseEncryptSecret: '',
      clientSecret: '',
    });
    expect(t).toThrow(/^database_encrypt_secret.*/);

    getOAuthSecretsSpy.mockReturnValue({
      sessionSecret: 'asdf',
      databaseEncryptSecret: 'asdf',
      clientSecret: '',
    });
    expect(t).toThrow(/^The length of.*/);

    getOAuthSecretsSpy.mockReturnValue({
      sessionSecret: 'asdf',
      databaseEncryptSecret: 'database_secret_'.repeat(2),
      clientSecret: '',
    });
    expect(t).not.toThrow();
  });

  it('should fail for missing Mongo values', () => {
    const config: Omit<BridgeOption, 'mongo'> & { mongo: Partial<MongoConfig> } = {
      api: {
        url: 'http://localhost',
        token: defaultAPIToken,
      },
      oauth: {
        enabled: true,
        discoveryURL: 'smth',
        clientID: 'id',
        baseURL: 'url',
      },
      mongo: {},
    };
    const t = (): BridgeConfiguration => getConfiguration(config);
    expect(t).toThrow(/Could not construct mongodb connection string.*/);
    config.mongo.host = 'mongo://';
    expect(t).toThrow(/Could not construct mongodb connection string.*/);
    config.mongo.password = 'pwd';
    expect(t).toThrow(/Could not construct mongodb connection string.*/);
    config.mongo.user = 'usr';
    expect(t).not.toThrow();
  });

  it('should set values using env var', () => {
    const config = getBasicEnvVar();
    config.LOGGING_COMPONENTS = 'a=true,b=false,c=true';
    const result = getConfiguration(envToConfiguration(config));
    expect(result.logging.destination).toBe(LogDestination.STDOUT);
    expect(result.logging.enabledComponents).toStrictEqual({
      a: true,
      b: false,
      c: true,
    });
  });

  it('should correctly eval booleans', () => {
    const config = envToConfiguration(getBasicEnvVar());
    const result = getConfiguration({
      ...config,
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

  it('should correctly map process.env to options', () => {
    const config: Record<EnvVar, string> = {
      LOGGING_COMPONENTS: 'xyz',
      SHOW_API_TOKEN: 'invalidBool',
      API_URL: 'apiUrl',
      API_TOKEN: 'apiToken',
      AUTH_MSG: 'authMsg',
      BASIC_AUTH_USERNAME: 'basicUsername',
      BASIC_AUTH_PASSWORD: 'basicPassword',
      REQUEST_TIME_LIMIT: '100',
      REQUESTS_WITHIN_TIME: '5',
      CLEAN_BUCKET_INTERVAL: '2000',
      OAUTH_ALLOWED_LOGOUT_URLS: 'logoutUrl',
      OAUTH_BASE_URL: 'baseUrl',
      OAUTH_CLIENT_ID: 'clientID',
      OAUTH_DISCOVERY: 'discovery',
      OAUTH_ENABLED: '',
      OAUTH_NAME_PROPERTY: 'name',
      OAUTH_SCOPE: 'scope',
      SECURE_COOKIE: 'fAlSe',
      SESSION_TIMEOUT_MIN: '2',
      TRUST_PROXY: '1',
      SESSION_VALIDATING_TIMEOUT_MIN: 'invalidNumber',
      OAUTH_ID_TOKEN_ALG: 'alg',
      CLI_DOWNLOAD_LINK: 'cliLink',
      INTEGRATIONS_PAGE_LINK: 'integrationLink',
      LOOK_AND_FEEL_URL: 'lookAndFeel',
      AUTOMATIC_PROVISIONING_MSG: 'automaticProvMsg',
      CONFIG_DIR: 'config',
      KEPTN_INSTALLATION_TYPE: 'installation',
      PROJECTS_PAGE_SIZE: '70',
      PREFIX_PATH: 'prefix',
      ENABLE_VERSION_CHECK: '0',
      MONGODB_DATABASE: 'mongoDatabase',
      MONGODB_HOST: 'mongoHost',
      MONGODB_PASSWORD: 'mongoPwd',
      MONGODB_USER: 'mongoUser',
      NODE_ENV: 'test',
      VERSION: '0.0.0',
    };

    const result = envToConfiguration(config);

    expect(result).toEqual({
      api: {
        showToken: true,
        token: 'apiToken',
        url: 'apiUrl',
      },
      auth: {
        authMessage: 'authMsg',
        basicPassword: 'basicPassword',
        basicUsername: 'basicUsername',
        cleanBucketIntervalMs: 2000,
        nRequestWithinTime: 5,
        requestTimeLimitMs: 100,
      },
      features: {
        automaticProvisioningMessage: 'automaticProvMsg',
        configDir: 'config',
        installationType: 'installation',
        pageSize: {
          project: 70,
          service: undefined,
        },
        prefixPath: 'prefix',
        versionCheck: false,
      },
      logging: {
        enabledComponents: 'xyz',
      },
      mode: 'test',
      mongo: {
        db: 'mongoDatabase',
        host: 'mongoHost',
        password: 'mongoPwd',
        user: 'mongoUser',
      },
      oauth: {
        allowedLogoutURL: 'logoutUrl',
        baseURL: 'baseUrl',
        clientID: 'clientID',
        discoveryURL: 'discovery',
        enabled: undefined,
        nameProperty: 'name',
        scope: 'scope',
        session: {
          secureCookie: false,
          timeoutMin: 2,
          trustProxyHops: 1,
          validationTimeoutMin: undefined,
        },
        tokenAlgorithm: 'alg',
      },
      urls: {
        CLI: 'cliLink',
        integrationPage: 'integrationLink',
        lookAndFeel: 'lookAndFeel',
      },
      version: '0.0.0',
    });
  });

  function getBasicEnvVar(): { [key in EnvVar]?: string } {
    return {
      [EnvVar.API_URL]: defaultAPIURL,
      [EnvVar.API_TOKEN]: defaultAPIToken,
      [EnvVar.MONGODB_HOST]: '',
      [EnvVar.MONGODB_USER]: '',
      [EnvVar.MONGODB_PASSWORD]: '',
      [EnvVar.NODE_ENV]: 'test',
    };
  }

  function getDefaultOAuthConfig(): { [key in EnvVar]?: string } {
    return {
      [EnvVar.API_URL]: defaultAPIURL,
      [EnvVar.API_TOKEN]: defaultAPIToken,
      [EnvVar.MONGODB_HOST]: 'asdf',
      [EnvVar.MONGODB_USER]: 'asdf',
      [EnvVar.MONGODB_PASSWORD]: 'asdf',
      [EnvVar.NODE_ENV]: 'test',
      [EnvVar.OAUTH_ENABLED]: 'true',
      [EnvVar.OAUTH_DISCOVERY]: 'asdf',
      [EnvVar.OAUTH_CLIENT_ID]: 'asdf',
      [EnvVar.OAUTH_BASE_URL]: 'asdf',
    };
  }
});
