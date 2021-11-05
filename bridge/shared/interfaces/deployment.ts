import { Sequence, SequenceState } from '../models/sequence';
import { Trace } from '../models/trace';
import { ResultTypes } from '../models/result-types';

export interface Deployment {
  stages: IStageDeployment[];
  keptnContext: string;
  service: string;
  image?: string;
  labels: { [key: string]: string };
  state: SequenceState; // useful for polling; if finished then just fetch/update openRemediations
}

export interface SubSequence {
  name: string;
  time: number;
  state: SequenceState;
  type: string;
  result: ResultTypes;
  id: string;
  message: string;
  hasPendingApproval: boolean;
}

export interface IStageDeployment {
  name: string;
  lastTimeUpdated: number;
  hasEvaluation: boolean;
  latestEvaluation?: Trace;
  openRemediations: Sequence[];
  remediationConfig?: string;
  approvalInformation?: {
    trace: Trace;
    deployedImage?: string;
  };
  subSequences: SubSequence[];
  deploymentURL?: string;
}
