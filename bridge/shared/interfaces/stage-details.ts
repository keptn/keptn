import { Approval } from './approval';
import { Trace } from '../models/trace';
import { Remediation } from '../models/remediation';

export interface StageDetails {
  services: StageServiceDetails[];
}

export interface StageServiceDetails {
  name: string;
  deployedURL?: string;
  deployedImage?: string;
  openApprovals?: Approval[]; // if undefined, there are no updates
  openRemediations?: Remediation[];
  failedEvaluation?: Trace;
}
