import { IClientFeatureFlags, IServerFeatureFlags } from '../shared/interfaces/feature-flags';

export class ClientFeatureFlags implements IClientFeatureFlags {}

export class ServerFeatureFlags implements IServerFeatureFlags {
  OAUTH_ENABLED = false;
}
