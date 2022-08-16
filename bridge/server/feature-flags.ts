import { IClientFeatureFlags, IServerFeatureFlags } from '../shared/interfaces/feature-flags';

export class ClientFeatureFlags implements IClientFeatureFlags {
  D3_ENABLED = true;
}

export class ServerFeatureFlags implements IServerFeatureFlags {
  OAUTH_ENABLED = false;
}
