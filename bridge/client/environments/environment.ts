import {DynamicEnvironment} from "./environment.dynamic";

class Environment extends DynamicEnvironment {

  public appTitle: string;
  public logoUrl: string;
  public logoInvertedUrl: string;
  public production: boolean;

  constructor() {
    super();
    this.production = false;
  }
}

export const environment = new Environment();
