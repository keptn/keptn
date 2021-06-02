declare var window: any;

export class DynamicEnvironment {
  public get config() {
    return window.config;
  }
}
