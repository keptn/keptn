import { KeptnInfoResult } from './keptn-info-result';
import { IMetadata } from './metadata';
import { KeptnVersions } from './keptn-versions';

export interface BridgeInfo {
  info: KeptnInfoResult;
  metadata: IMetadata | null;
  versions?: KeptnVersions;
}
