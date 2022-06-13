import { EnabledComponents, LogDestination } from '../utils/logger';

/**
 * Option to configure the Bridge-Server
 */
export interface BridgeOption {
  logging?: LogOptions;
  auth?: AuthOptions;
  oauth?: OAuthOptions;
  api?: APIOptions;
  feature?: FeatureOption;
  version?: string;
  mode?: string;
  mongo?: MongoOptions;
}

export enum EnvType {
  PRODUCTION = 'production',
  TEST = 'test',
  DEV = 'development',
}

interface LogOptions {
  destination?: LogDestination;
  enabledComponents?: string;
}

interface AuthOptions {
  authMessage?: string;
}

interface OAuthOptions {
  enabled?: boolean;
  discoveryURL?: string;
  baseURL?: string;
  clientID?: string;
}

interface APIOptions {
  url?: string;
  token?: string;
  showToken?: boolean;
}

interface MongoOptions {
  user: string;
  password: string;
  host: string;
}

interface FeatureOption {
  automaticProvisioningMessage?: string;
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
  clientSecret?: string;
  scope: string;
  tokenAlgorithm: string;
  allowedLogoutURL: string;
  nameProperty?: string;
  session: OAuthSessionConfig;
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
  integrationPage: string;
  CLI: string;
}

export interface MongoConfig {
  user: string;
  password: string;
  host: string;
  db: string;
}
