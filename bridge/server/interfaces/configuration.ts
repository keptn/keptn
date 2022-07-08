import { EnabledComponents, LogDestination } from '../utils/logger';

export interface OAuthSecrets {
  sessionSecret: string;
  databaseEncryptSecret: string;
  clientSecret?: string;
}

type RecursivePartial<T> = {
  [P in keyof T]?: RecursivePartial<T[P]>;
};
/**
 * Option to configure the Bridge-Server
 */
export type BridgeOption = Omit<RecursivePartial<BridgeConfiguration>, 'logging' | 'mode'> & {
  logging?: LogOptions;
  mode?: string;
};

export enum EnvType {
  PRODUCTION = 'production',
  TEST = 'test',
  DEV = 'development',
}

interface LogOptions {
  destination?: LogDestination;
  enabledComponents?: string;
}

/**
 * Configuration object
 */
export interface BridgeConfiguration {
  logging: LogConfiguration;
  api: APIConfig;
  auth: AuthConfig;
  oauth: OAuthConfig;
  urls: URLsConfig;
  features: FeatureConfig;
  version: string;
  mode: EnvType;
  mongo: MongoConfig;
}

export interface LogConfiguration {
  destination: LogDestination;
  enabledComponents: EnabledComponents;
}

export interface AuthConfig {
  requestTimeLimitMs: number;
  nRequestWithinTime: number;
  cleanBucketIntervalMs: number;
  basicUsername?: string;
  basicPassword?: string;
  authMessage: string;
}

export interface OAuthConfig {
  enabled: boolean;
  discoveryURL: string;
  baseURL: string;
  clientID: string;
  scope: string;
  tokenAlgorithm: string;
  allowedLogoutURL: string;
  nameProperty?: string;
  session: OAuthSessionConfig;
  secrets: OAuthSecrets;
}

export interface OAuthSessionConfig {
  secureCookie: boolean;
  trustProxyHops: number;
  timeoutMin: number;
  validationTimeoutMin: number;
}

export interface APIConfig {
  url: string;
  token: string | undefined;
  showToken: boolean;
}

export interface FeatureConfig {
  pageSize: PageSizeConfiguration;
  installationType: string;
  automaticProvisioningMessage: string;
  prefixPath: string;
  configDir: string;
  versionCheck: boolean;
}

export interface PageSizeConfiguration {
  project: number;
  service: number;
}

export interface URLsConfig {
  lookAndFeel?: string;
  CLI: string;
}

export interface MongoConfig {
  user: string;
  password: string;
  host: string;
  db: string;
  externalConnectionString?: string;
}

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
