import { KeptnInfoResult } from './keptn-info-result';
import { KeptnVersions } from './keptn-versions';
import { Metadata } from './metadata';

export interface KeptnInfo {
  bridgeInfo: KeptnInfoResult;
  availableVersions?: KeptnVersions;
  versionCheckEnabled?: boolean;
  metadata: Metadata;
  authCommand?: string;
  keptnVersionInvalid?: boolean;
  bridgeVersionInvalid?: boolean;
}
