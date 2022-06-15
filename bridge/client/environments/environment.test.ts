import { DynamicEnvironment } from './environment.dynamic';

class Environment extends DynamicEnvironment {
  constructor() {
    super();
    this.pollingIntervalMillis = 0;
  }
}

export const environment = new Environment();
