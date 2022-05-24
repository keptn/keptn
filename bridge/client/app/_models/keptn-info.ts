import { KeptnInfoResult } from '../../../shared/interfaces/keptn-info-result';
import { KeptnVersions } from './keptn-versions';

export interface KeptnInfo {
  bridgeInfo: KeptnInfoResult;
  availableVersions?: KeptnVersions;
  versionCheckEnabled?: boolean;
  authCommand?: string;
  keptnVersionInvalid?: boolean;
  bridgeVersionInvalid?: boolean;
}
