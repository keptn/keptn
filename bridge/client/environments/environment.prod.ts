import { DynamicEnvironment } from './environment.dynamic';

class Environment extends DynamicEnvironment {
  public appTitle?: string;
  public logoUrl?: string;
  public logoInvertedUrl?: string;
  public production: boolean;
  public appConfigUrl: string;
  public baseUrl: string;

  constructor() {
    super();
    this.production = true;
    this.appConfigUrl = 'assets/branding/app-config.json';
    this.baseUrl = '/bridge';
  }
}

export const environment = new Environment();
