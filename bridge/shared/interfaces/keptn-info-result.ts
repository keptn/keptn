import { IClientFeatureFlags } from './feature-flags';

export interface KeptnInfoResult {
  featureFlags: IClientFeatureFlags;
  bridgeVersion?: string;
  apiUrl?: string;
  apiToken?: string;
  cliDownloadLink: string;
  enableVersionCheckFeature: boolean;
  showApiToken: boolean;
  projectsPageSize?: number;
  servicesPageSize?: number;
  authType: string;
  user?: string;
}
