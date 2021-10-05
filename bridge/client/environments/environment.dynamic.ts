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
