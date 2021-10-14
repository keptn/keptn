import { FeatureFlags } from '../../shared/interfaces/feature-flags';

export interface WindowConfig {
  appTitle: string;
  logoUrl: string;
  logoInvertedUrl: string;
  stylesheetUrl?: string;
}

declare global {
  interface Window {
    config: WindowConfig;
  }
}

export class DynamicEnvironment {
  public appTitle?: string;
  public logoUrl?: string;
  public logoInvertedUrl?: string;
  public production: boolean;
  public appConfigUrl: string;
  public baseUrl: string;
  public featureFlags: FeatureFlags;
  public pollingIntervalMillis?: number;

  constructor() {
    this.production = false;
    this.appConfigUrl = 'assets/default-branding/app-config.json';
    this.baseUrl = '/';
    this.featureFlags = {
      exampleFlag: true,
    };
  }

  public get config(): WindowConfig {
    return (
      window.config || {
        appTitle: 'keptn',
        logoUrl: 'assets/default-branding/logo.png',
        logoInvertedUrl: 'assets/default-branding/logo_inverted.png',
      }
    );
  }
}
