import { Sequence } from '../models/sequence';

export interface StageOverviewServiceInfo {
  latestSequence?: Sequence;
  name: string;
  deployedImage?: string;
}

export interface StageOverview {
  name: string;
  services: StageOverviewServiceInfo[];
  failedEvaluations?: number; // if undefined; no update
  openRemediations?: number;
  openApprovals?: number;
}
