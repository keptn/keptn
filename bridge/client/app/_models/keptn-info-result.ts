export interface KeptnInfoResult {
  bridgeVersion?: string;
  keptnInstallationType?: string;
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
