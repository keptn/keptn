import { DynamicEnvironment } from './environment.dynamic';

class Environment extends DynamicEnvironment {
  constructor() {
    super();
    this.pollingIntervalMillis = 0;
    this.featureFlags = {
      exampleFlag: false,
    };
  }
}

export const environment = new Environment();
