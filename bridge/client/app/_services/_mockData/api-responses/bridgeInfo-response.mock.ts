import { IClientFeatureFlags } from '../../../../../shared/interfaces/feature-flags';

const featureFlags: IClientFeatureFlags = {
  RESOURCE_SERVICE_ENABLED: false,
  D3_HEATMAP_ENABLED: false,
};

const bridgeInfo = {
  bridgeVersion: '0.10.0',
  featureFlags,
  keptnInstallationType: 'QUALITY_GATES,CONTINUOUS_OPERATIONS,CONTINUOUS_DELIVERY',
  apiUrl: 'http://example.com/api',
  apiToken: 'random_api_token',
  cliDownloadLink: 'https://github.com/keptn/keptn/releases/tag/0.10.0',
  enableVersionCheckFeature: true,
  showApiToken: true,
  authType: 'BASIC',
};
export { bridgeInfo as BridgeInfoResponseMock };
