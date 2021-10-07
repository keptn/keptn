import { DynamicEnvironment } from './environment.dynamic';

class Environment extends DynamicEnvironment {
  public appTitle?: string;
  public logoUrl?: string;
  public logoInvertedUrl?: string;
  public production: boolean;
  public appConfigUrl: string;
  public baseUrl: string;
  public versionCheckEnabled: boolean;

  constructor() {
    super();
    this.production = false;
    this.appConfigUrl = 'assets/default-branding/app-config.json';
    this.baseUrl = '/';
    this.versionCheckEnabled = false;
  }
}

export const environment = new Environment();
