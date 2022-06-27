import { DynamicEnvironment } from './environment.dynamic';

class Environment extends DynamicEnvironment {
  constructor() {
    super();
    this.pollingIntervalMillis = 0;
    this.production = true;
  }
}

export const environment = new Environment();
