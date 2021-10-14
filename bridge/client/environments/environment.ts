import { DynamicEnvironment } from './environment.dynamic';

class Environment extends DynamicEnvironment {
  constructor() {
    super();
    this.featureFlags = {
      exampleFlag: true,
    };
  }
}

export const environment = new Environment();
