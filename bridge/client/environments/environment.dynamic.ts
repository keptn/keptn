export interface WindowConfig {
  appTitle: string;
  logoUrl: string;
  logoInvertedUrl: string;
  stylesheetUrl?: string;
}

export class DynamicEnvironment {
  public production: boolean;
  public appConfigUrl: string;
  public baseUrl: string;
  public pollingIntervalMillis?: number;
  public ngrx: boolean;
  public config: WindowConfig = {
    appTitle: 'Keptn',
    logoUrl: 'assets/default-branding/logo.png',
    logoInvertedUrl: 'assets/default-branding/logo_inverted.png',
  };

  constructor() {
    this.production = false;
    this.appConfigUrl = 'assets/default-branding/app-config.json';
    this.baseUrl = '/';
    this.ngrx = true;
  }
}
