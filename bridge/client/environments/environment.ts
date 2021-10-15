import { DynamicEnvironment } from './environment.dynamic';

class Environment extends DynamicEnvironment {
  constructor() {
    super();
    this.featureFlags = {
    };
  }
}

export const environment = new Environment();
