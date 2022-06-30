import { execSync } from 'child_process';
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';
import {
  APIConfig,
  AuthConfig,
  BridgeConfiguration,
  BridgeOption,
  EnvType,
  FeatureConfig,
  LogConfiguration,
  MongoConfig,
  OAuthConfig,
  URLsConfig,
} from '../interfaces/configuration';

import { EnabledComponents, LogDestination, logger as log } from './logger';

/**
 * Env var list
 */

export enum EnvVar {
  LOGGING_COMPONENTS = 'LOGGING_COMPONENTS',
  SHOW_API_TOKEN = 'SHOW_API_TOKEN',
  API_URL = 'API_URL',
  API_TOKEN = 'API_TOKEN',
  AUTH_MSG = 'AUTH_MSG',
  BASIC_AUTH_USERNAME = 'BASIC_AUTH_USERNAME',
  BASIC_AUTH_PASSWORD = 'BASIC_AUTH_PASSWORD',
  REQUEST_TIME_LIMIT = 'REQUEST_TIME_LIMIT',
  REQUESTS_WITHIN_TIME = 'REQUESTS_WITHIN_TIME',
  CLEAN_BUCKET_INTERVAL = 'CLEAN_BUCKET_INTERVAL',
  OAUTH_ALLOWED_LOGOUT_URLS = 'OAUTH_ALLOWED_LOGOUT_URLS',
  OAUTH_BASE_URL = 'OAUTH_BASE_URL',
  OAUTH_CLIENT_ID = 'OAUTH_CLIENT_ID',
  OAUTH_CLIENT_SECRET = 'OAUTH_CLIENT_SECRET',
  OAUTH_DISCOVERY = 'OAUTH_DISCOVERY',
  OAUTH_ENABLED = 'OAUTH_ENABLED',
  OAUTH_NAME_PROPERTY = 'OAUTH_NAME_PROPERTY',
  OAUTH_SCOPE = 'OAUTH_SCOPE',
  SECURE_COOKIE = 'SECURE_COOKIE',
  SESSION_TIMEOUT_MIN = 'SESSION_TIMEOUT_MIN',
  TRUST_PROXY = 'TRUST_PROXY',
  SESSION_VALIDATING_TIMEOUT_MIN = 'SESSION_VALIDATING_TIMEOUT_MIN',
  OAUTH_ID_TOKEN_ALG = 'OAUTH_ID_TOKEN_ALG',
  CLI_DOWNLOAD_LINK = 'CLI_DOWNLOAD_LINK',
  INTEGRATIONS_PAGE_LINK = 'INTEGRATIONS_PAGE_LINK',
  LOOK_AND_FEEL_URL = 'LOOK_AND_FEEL_URL',
  AUTOMATIC_PROVISIONING_MSG = 'AUTOMATIC_PROVISIONING_MSG',
  CONFIG_DIR = 'CONFIG_DIR',
  KEPTN_INSTALLATION_TYPE = 'KEPTN_INSTALLATION_TYPE',
  PROJECTS_PAGE_SIZE = 'PROJECTS_PAGE_SIZE',
  SERVICES_PAGE_SIZE = 'SERVICES_PAGE_SIZE',
  PREFIX_PATH = 'PREFIX_PATH',
  ENABLE_VERSION_CHECK = 'ENABLE_VERSION_CHECK',
  MONGODB_DATABASE = 'MONGODB_DATABASE',
  MONGODB_HOST = 'MONGODB_HOST',
  MONGODB_PASSWORD = 'MONGODB_PASSWORD',
  MONGODB_USER = 'MONGODB_USER',
  NODE_ENV = 'NODE_ENV',
  VERSION = 'VERSION',
}

const _componentName = 'Configuration';

/**
 * @param options Customization options that override env var options.
 * @returns Returns the Bridge-server configuration.
 */
export function getConfiguration(options?: BridgeOption): BridgeConfiguration {
  const logConfig = getLogConfiguration(options);
  const apiConfig = getAPIConfiguration(options);
  const authConfig = getAuthConfiguration(apiConfig, options);
  const oauthConfig = getOAuthConfiguration(options);
  const urlsConfig = getURLsConfiguration();
  const featConfig = getFeatureConfiguration(options);
  const mongoConfig = getMongoConfiguration(options, oauthConfig.enabled);

  // mode and version
  const _mode = options?.mode ?? (process.env[EnvVar.NODE_ENV] || 'development');
  const modeMap: Record<string, EnvType> = {
    production: EnvType.PRODUCTION,
    test: EnvType.TEST,
    development: EnvType.DEV,
  };
  const mode = modeMap[_mode];
  const version = options?.version ?? process.env[EnvVar.VERSION] ?? 'develop';

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

function getLogConfiguration(options?: BridgeOption): LogConfiguration {
  const logDestination = options?.logging?.destination ?? LogDestination.STDOUT;
  const loggingComponents = Object.create({}) as EnabledComponents;
  const loggingComponentsString = options?.logging?.enabledComponents ?? process.env[EnvVar.LOGGING_COMPONENTS] ?? '';
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

function getAPIConfiguration(options?: BridgeOption): APIConfig {
  const _showToken = options?.api?.showToken ?? toBool(process.env[EnvVar.SHOW_API_TOKEN] ?? 'true');
  const apiUrl = options?.api?.url ?? process.env[EnvVar.API_URL];
  if (!apiUrl) {
    throw new Error('API_URL is not provided');
  }
  const apiToken = options?.api?.token ?? process.env[EnvVar.API_TOKEN];
  if (typeof apiToken !== 'string') {
    log.warning(_componentName, 'API_TOKEN was not provided.');
  }

  return {
    showToken: _showToken,
    url: apiUrl,
    token: apiToken,
  };
}

function getAuthConfiguration(api: APIConfig, options?: BridgeOption): AuthConfig {
  const authMsg =
    (options?.auth?.authMessage ?? process.env[EnvVar.AUTH_MSG]) ||
    `keptn auth --endpoint=${api.url} --api-token=${api.token}`;
  const basicUser = process.env[EnvVar.BASIC_AUTH_USERNAME];
  const basicPass = process.env[EnvVar.BASIC_AUTH_PASSWORD];
  const requestLimit = toInt(process.env[EnvVar.REQUEST_TIME_LIMIT], 60) * 60 * 1000;
  const requestWithinTime = toInt(process.env[EnvVar.REQUESTS_WITHIN_TIME], 10);
  const cleanBucket = toInt(process.env[EnvVar.CLEAN_BUCKET_INTERVAL], 60) * 60 * 1000;

  return {
    authMessage: authMsg,
    basicUsername: basicUser,
    basicPassword: basicPass,
    requestTimeLimitMs: requestLimit,
    nRequestWithinTime: requestWithinTime,
    cleanBucketIntervalMs: cleanBucket,
  };
}

function getOAuthConfiguration(options?: BridgeOption): OAuthConfig {
  const logoutURL = process.env[EnvVar.OAUTH_ALLOWED_LOGOUT_URLS] ?? '';
  let baseURL = options?.oauth?.baseURL ?? process.env[EnvVar.OAUTH_BASE_URL];
  let clientID = options?.oauth?.clientID ?? process.env[EnvVar.OAUTH_CLIENT_ID];
  const clientSecret = process.env[EnvVar.OAUTH_CLIENT_SECRET];
  let discoveryURL = options?.oauth?.discoveryURL ?? process.env[EnvVar.OAUTH_DISCOVERY];
  const enabled = options?.oauth?.enabled ?? toBool(process.env[EnvVar.OAUTH_ENABLED] ?? 'false');
  const nameProperty = process.env[EnvVar.OAUTH_NAME_PROPERTY];
  const scope = process.env[EnvVar.OAUTH_SCOPE] ?? '';
  const secureCookie = toBool(process.env[EnvVar.SECURE_COOKIE] ?? 'false');
  const timeout = toInt(process.env[EnvVar.SESSION_TIMEOUT_MIN], 60);
  const proxyHops = toInt(process.env[EnvVar.TRUST_PROXY], 1);
  const validation = toInt(process.env[EnvVar.SESSION_VALIDATING_TIMEOUT_MIN], 60);
  const algo = process.env[EnvVar.OAUTH_ID_TOKEN_ALG] || 'RS256';
  if (enabled) {
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
  } else {
    discoveryURL = '';
    clientID = '';
    baseURL = '';
  }

  return {
    allowedLogoutURL: logoutURL,
    baseURL: baseURL,
    clientID: clientID,
    clientSecret: clientSecret,
    discoveryURL: discoveryURL,
    enabled: enabled,
    nameProperty: nameProperty,
    scope: scope.trim(),
    session: {
      secureCookie: secureCookie,
      timeoutMin: timeout,
      trustProxyHops: proxyHops,
      validationTimeoutMin: validation,
    },
    tokenAlgorithm: algo,
  };
}

function getURLsConfiguration(): URLsConfig {
  const cliURL = process.env[EnvVar.CLI_DOWNLOAD_LINK] ?? 'https://github.com/keptn/keptn/releases';
  const integrationURL = process.env[EnvVar.INTEGRATIONS_PAGE_LINK] ?? 'https://get.keptn.sh/integrations.html';
  const looksURL = process.env[EnvVar.LOOK_AND_FEEL_URL];

  return {
    CLI: cliURL,
    integrationPage: integrationURL,
    lookAndFeel: looksURL,
  };
}

function getFeatureConfiguration(options?: BridgeOption): FeatureConfig {
  const provisioningMsg =
    options?.feature?.automaticProvisioningMessage ?? process.env[EnvVar.AUTOMATIC_PROVISIONING_MSG] ?? '';
  const configDir =
    process.env[EnvVar.CONFIG_DIR] ?? join(dirname(fileURLToPath(import.meta.url)), '../../../../config');
  const installationType =
    process.env[EnvVar.KEPTN_INSTALLATION_TYPE] ?? 'QUALITY_GATES,CONTINUOUS_OPERATIONS,CONTINUOUS_DELIVERY';
  const projectSize = toInt(process.env[EnvVar.PROJECTS_PAGE_SIZE], 50); // client\app\_services\api.service.ts
  const serviceSize = toInt(process.env[EnvVar.SERVICES_PAGE_SIZE], 50); // no use
  const prefixPath = process.env[EnvVar.PREFIX_PATH] ?? '/';
  const versionCheck = toBool(process.env[EnvVar.ENABLE_VERSION_CHECK] ?? 'true');

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

function getMongoConfiguration(options?: BridgeOption, doChecks = true): MongoConfig {
  const db = process.env[EnvVar.MONGODB_DATABASE] ?? 'openid';
  const host = options?.mongo?.host ?? process.env[EnvVar.MONGODB_HOST];
  const pwd = options?.mongo?.password ?? process.env[EnvVar.MONGODB_PASSWORD];
  const usr = options?.mongo?.user ?? process.env[EnvVar.MONGODB_USER];
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
function toInt(v: string | undefined, d: number): number {
  if (v) {
    const val = parseInt(v, 10);
    if (!isNaN(val)) {
      return val;
    }
  }
  return d;
}
