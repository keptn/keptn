declare var window: any;

export class DynamicEnvironment {
  public get config() {
    return window.config || {
      "appTitle": "keptn",
      "logoUrl": "assets/default-branding/logo.png",
      "logoInvertedUrl": "assets/default-branding/logo_inverted.png"
    };
  }
}
