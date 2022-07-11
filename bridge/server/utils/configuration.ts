import { dirname, join } from 'path';
import { fileURLToPath } from 'url';
import {
  APIConfig,
  AuthConfig,
  BridgeConfiguration,
  BridgeOption,
  EnvType,
  EnvVar,
  FeatureConfig,
  LogConfiguration,
  MongoConfig,
  OAuthConfig,
  OAuthSecrets,
  URLsConfig,
} from '../interfaces/configuration';

import { EnabledComponents, LogDestination, logger as log } from './logger';
import { getOAuthMongoExternalConnectionString, getOAuthSecrets } from '../user/secrets';

const _componentName = 'Configuration';

/**
 * @param options Customization options that override env var options.
 * @returns Returns the Bridge-server configuration.
 */
export function getConfiguration(options: BridgeOption): BridgeConfiguration {
  const featConfig = getFeatureConfiguration(options);
  const logConfig = getLogConfiguration(options);
  const apiConfig = getAPIConfiguration(options);
  const authConfig = getAuthConfiguration(apiConfig, options);
  const oauthConfig = getOAuthConfiguration(featConfig.configDir, options);
  const urlsConfig = getURLsConfiguration(options);
  const mongoConfig = getMongoConfiguration(featConfig.configDir, options, oauthConfig.enabled);

  // mode and version
  const _mode = options.mode || 'development';
  const modeMap: Record<string, EnvType> = {
    production: EnvType.PRODUCTION,
    test: EnvType.TEST,
    development: EnvType.DEV,
  };
  const mode = modeMap[_mode];
  const version = options.version ?? 'develop';

  return {
    logging: logConfig,
    api: apiConfig,
    auth: authConfig,
    oauth: oauthConfig,
    urls: urlsConfig,
    features: featConfig,
    mode: mode,
    mongo: mongoConfig,
    version: version,
  };
}

function getLogConfiguration(options: BridgeOption): LogConfiguration {
  const logDestination = options.logging?.destination ?? LogDestination.STDOUT;
  const loggingComponents = Object.create({}) as EnabledComponents;
  const loggingComponentsString = options.logging?.enabledComponents ?? '';
  if (loggingComponentsString.length > 0) {
    const components = loggingComponentsString.split(',').map((s) => s.trim());
    for (const component of components) {
      const [name, value] = parseComponent(component);
      loggingComponents[name] = value;
    }
  }
  return {
    destination: logDestination,
    enabledComponents: loggingComponents,
  };
}

function getAPIConfiguration(options: BridgeOption): APIConfig {
  const _showToken = options.api?.showToken ?? true;
  const apiUrl = options.api?.url;
  const apiToken = options.api?.token;

  if (!apiUrl) {
    throw new Error('API_URL is not provided');
  }

  if (typeof apiToken !== 'string') {
    log.warning(_componentName, 'API_TOKEN was not provided.');
  }

  return {
    showToken: _showToken,
    url: apiUrl,
    token: apiToken,
  };
}

function getAuthConfiguration(api: APIConfig, options: BridgeOption): AuthConfig {
  const authMsg = options.auth?.authMessage || `keptn auth --endpoint=${api.url} --api-token=${api.token}`;
  const basicUser = options.auth?.basicUsername;
  const basicPass = options.auth?.basicPassword;
  const requestLimit = (options.auth?.requestTimeLimitMs ?? 60) * 60 * 1000;
  const requestWithinTime = options.auth?.nRequestWithinTime ?? 10;
  const cleanBucket = (options.auth?.cleanBucketIntervalMs ?? 60) * 60 * 1000;

  return {
    authMessage: authMsg,
    basicUsername: basicUser,
    basicPassword: basicPass,
    requestTimeLimitMs: requestLimit,
    nRequestWithinTime: requestWithinTime,
    cleanBucketIntervalMs: cleanBucket,
  };
}

function getOAuthConfiguration(configDir: string, options: BridgeOption): OAuthConfig {
  const logoutURL = options.oauth?.allowedLogoutURL ?? '';
  const enabled = options.oauth?.enabled ?? false;
  const nameProperty = options.oauth?.nameProperty;
  const scope = options.oauth?.scope ?? '';
  const secureCookie = options.oauth?.session?.secureCookie ?? false;
  const timeout = options.oauth?.session?.timeoutMin ?? 60;
  const proxyHops = options.oauth?.session?.trustProxyHops ?? 1;
  const validation = options.oauth?.session?.validationTimeoutMin ?? 60;
  const algo = options.oauth?.tokenAlgorithm || 'RS256';
  const basicOptions = validateOAuthBasicOptions(
    enabled,
    options.oauth?.discoveryURL,
    options.oauth?.clientID,
    options.oauth?.baseURL,
    configDir
  );

  return {
    allowedLogoutURL: logoutURL,
    baseURL: basicOptions.baseURL,
    clientID: basicOptions.clientID,
    discoveryURL: basicOptions.discoveryURL,
    enabled: enabled,
    nameProperty: nameProperty,
    scope: scope.trim(),
    session: {
      secureCookie: secureCookie,
      timeoutMin: timeout,
      trustProxyHops: proxyHops,
      validationTimeoutMin: validation,
    },
    secrets: basicOptions.secrets ?? {
      sessionSecret: '',
      databaseEncryptSecret: '',
    },
    tokenAlgorithm: algo,
  };
}

function validateOAuthBasicOptions(
  enabled: boolean,
  discoveryURL: string | undefined,
  clientID: string | undefined,
  baseURL: string | undefined,
  configDir: string
): { discoveryURL: string; clientID: string; baseURL: string; secrets?: OAuthSecrets } {
  if (enabled) {
    const secrets = getOAuthSecrets(configDir);
    const errorSuffix =
      'must be defined when OAuth based login (OAUTH_ENABLED) is activated.' +
      ' Please check your environment variables.';
    if (!discoveryURL) {
      throw new Error(`OAUTH_DISCOVERY ${errorSuffix}`);
    }
    if (!clientID) {
      throw new Error(`OAUTH_CLIENT_ID ${errorSuffix}`);
    }
    if (!baseURL) {
      throw new Error(`OAUTH_BASE_URL ${errorSuffix}`);
    }
    validateSecrets(secrets);
    return {
      discoveryURL,
      baseURL,
      clientID,
      secrets,
    };
  }
  return {
    discoveryURL: '',
    clientID: '',
    baseURL: '',
  };
}

function validateSecrets(secrets: OAuthSecrets): void {
  const errorSuffix =
    'must be defined when OAuth based login (OAUTH_ENABLED) is activated. Please check your bridge-oauth secret.';
  if (!secrets.sessionSecret) {
    throw Error(`session_secret ${errorSuffix}`);
  }

  if (!secrets.databaseEncryptSecret) {
    throw Error(`database_encrypt_secret ${errorSuffix}`);
  } else if (secrets.databaseEncryptSecret.length !== 32) {
    throw Error(`The length of the secret "database_encrypt_secret" must be 32`);
  }
}

function getURLsConfiguration(options: BridgeOption): URLsConfig {
  const cliURL = options.urls?.CLI ?? 'https://github.com/keptn/keptn/releases';
  const integrationURL = options.urls?.integrationPage ?? 'https://get.keptn.sh/integrations.html';
  const looksURL = options.urls?.lookAndFeel;

  return {
    CLI: cliURL,
    integrationPage: integrationURL,
    lookAndFeel: looksURL,
  };
}

function getFeatureConfiguration(options: BridgeOption): FeatureConfig {
  const provisioningMsg = options.features?.automaticProvisioningMessage ?? '';
  const configDir = options.features?.configDir ?? join(dirname(fileURLToPath(import.meta.url)), '../../../../config');
  const installationType =
    options.features?.installationType ?? 'QUALITY_GATES,CONTINUOUS_OPERATIONS,CONTINUOUS_DELIVERY';
  const projectSize = options.features?.pageSize?.project ?? 50; // client\app\_services\api.service.ts
  const serviceSize = options.features?.pageSize?.service ?? 50; // no use
  const prefixPath = options.features?.prefixPath ?? '/';
  const versionCheck = options.features?.versionCheck ?? true;

  return {
    automaticProvisioningMessage: provisioningMsg,
    configDir: configDir,
    installationType: installationType,
    pageSize: {
      project: projectSize,
      service: serviceSize,
    },
    prefixPath: prefixPath,
    versionCheck: versionCheck,
  };
}

function getMongoConfiguration(configDir: string, options: BridgeOption, doChecks = true): MongoConfig {
  const db = options.mongo?.db ?? 'openid';
  const host = options.mongo?.host;
  const pwd = options.mongo?.password;
  const usr = options.mongo?.user;
  if (doChecks) {
    const errMsg =
      'Could not construct mongodb connection string: env vars "MONGODB_HOST", "MONGODB_USER" and "MONGODB_PASSWORD" have to be set';
    if (!host) {
      throw Error(errMsg);
    }
    if (typeof pwd !== 'string') {
      throw Error(errMsg);
    }
    if (typeof usr !== 'string') {
      throw Error(errMsg);
    }
  }

  return {
    db: db,
    host: host || '',
    password: pwd || '',
    user: usr || '',
    externalConnectionString: getOAuthMongoExternalConnectionString(configDir),
  };
}

function parseComponent(component: string): [string, boolean] {
  // we expect only componentName = bool
  const split = component.split('=', 3);
  return [split[0].trim(), toBool(split[1])];
}

/**
 * Convert string to boolean. If the input is equal to false or 0, it returns false. True otherwise.
 * @param v string to convert.
 */
function toBool(v: string): boolean {
  const val = v.toLowerCase();
  return val !== '0' && val !== 'false';
}

/**
 * Convert string to int. If the input cannot be converted, returns the default.
 * @param v string to convert.
 * @param d default value.
 */
function toInt(v: string, d = 0): number {
  if (v) {
    const val = parseInt(v, 10);
    if (!isNaN(val)) {
      return val;
    }
  }
  return d;
}

function isInt(v: string | undefined): v is string {
  if (!v) {
    return false;
  }
  const val = parseInt(v, 10);
  return !isNaN(val);
}

export function envToConfiguration(env: { [key in EnvVar]?: string }): BridgeOption {
  return {
    version: env.VERSION,
    mode: env.NODE_ENV,
    logging: {
      enabledComponents: env.LOGGING_COMPONENTS,
    },
    api: {
      showToken: env.SHOW_API_TOKEN ? toBool(env.SHOW_API_TOKEN) : undefined,
      url: env.API_URL,
      token: env.API_TOKEN,
    },
    auth: {
      authMessage: env.AUTH_MSG,
      basicPassword: env.BASIC_AUTH_PASSWORD,
      basicUsername: env.BASIC_AUTH_USERNAME,
      requestTimeLimitMs: isInt(env.REQUEST_TIME_LIMIT) ? toInt(env.REQUEST_TIME_LIMIT) : undefined,
      nRequestWithinTime: isInt(env.REQUESTS_WITHIN_TIME) ? toInt(env.REQUESTS_WITHIN_TIME) : undefined,
      cleanBucketIntervalMs: isInt(env.CLEAN_BUCKET_INTERVAL) ? toInt(env.CLEAN_BUCKET_INTERVAL) : undefined,
    },
    urls: {
      CLI: env.CLI_DOWNLOAD_LINK,
      integrationPage: env.INTEGRATIONS_PAGE_LINK,
      lookAndFeel: env.LOOK_AND_FEEL_URL,
    },
    features: {
      configDir: env.CONFIG_DIR,
      installationType: env.KEPTN_INSTALLATION_TYPE,
      automaticProvisioningMessage: env.AUTOMATIC_PROVISIONING_MSG,
      versionCheck: env.ENABLE_VERSION_CHECK ? toBool(env.ENABLE_VERSION_CHECK) : undefined,
      prefixPath: env.PREFIX_PATH,
      pageSize: {
        project: isInt(env.PROJECTS_PAGE_SIZE) ? toInt(env.PROJECTS_PAGE_SIZE) : undefined,
      },
    },
    mongo: {
      db: env.MONGODB_DATABASE,
      host: env.MONGODB_HOST,
      password: env.MONGODB_PASSWORD,
      user: env.MONGODB_USER,
    },
    oauth: {
      enabled: env.OAUTH_ENABLED ? toBool(env.OAUTH_ENABLED) : undefined,
      discoveryURL: env.OAUTH_DISCOVERY,
      nameProperty: env.OAUTH_NAME_PROPERTY,
      scope: env.OAUTH_SCOPE,
      clientID: env.OAUTH_CLIENT_ID,
      baseURL: env.OAUTH_BASE_URL,
      allowedLogoutURL: env.OAUTH_ALLOWED_LOGOUT_URLS,
      tokenAlgorithm: env.OAUTH_ID_TOKEN_ALG || 'RS256',
      session: {
        secureCookie: env.SECURE_COOKIE ? toBool(env.SECURE_COOKIE) : undefined,
        timeoutMin: isInt(env.SESSION_TIMEOUT_MIN) ? toInt(env.SESSION_TIMEOUT_MIN) : undefined,
        trustProxyHops: isInt(env.TRUST_PROXY) ? toInt(env.TRUST_PROXY) : undefined,
        validationTimeoutMin: isInt(env.SESSION_VALIDATING_TIMEOUT_MIN)
          ? toInt(env.SESSION_VALIDATING_TIMEOUT_MIN)
          : undefined,
      },
    },
  };
}
