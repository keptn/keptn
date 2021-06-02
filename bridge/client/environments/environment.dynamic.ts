declare var window: any;

export class DynamicEnvironment {
  public get config() {
    return window.config || {
      "appTitle": "keptn",
      "logoUrl": "assets/branding/logo.png",
      "logoInvertedUrl": "assets/branding/logo_inverted.png"
    };
  }
}
