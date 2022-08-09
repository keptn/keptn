import { IShipyardSequence } from './shipyard';
import { ResultTypes } from '../models/result-types';
import { StatusTypes } from '../models/status-types';

export interface SequenceExecutionResult {
  nextPageKey?: number;
  pageSize?: number;
  totalCount?: number;
  sequenceExecutions: SequenceExecution[];
}

export interface SequenceExecution {
  _id: string;
  inputProperties?: Record<string, unknown>;
  schemaVersion: string;
  scope: EventScope;
  sequence: IShipyardSequence;
  status: SequenceExecutionStatus;
  triggeredAt: string;
}

export interface EventScope {
  eventType: string;
  gitcommitid: string;
  keptnContext: string;
  labels: Record<string, string | undefined>;
  message: string;
  project: string;
  result: ResultTypes;
  service: string;
  stage: string;
  status: StatusTypes;
  triggeredId: string;
}

export interface SequenceExecutionStatus {
  currentTask: {
    events: [
      {
        eventType: string;
        properties: Record<string, unknown>;
        result: ResultTypes;
        source: string;
        status: StatusTypes;
        time: string; // ISO_INSTANT format e.g. yyyy-mm-ddThh:mm:ss.SSSZ
      }
    ];
    name: string;
    triggeredID: string;
  };
  previousTasks: [
    {
      name: string;
      properties: Record<string, unknown>;
      result: ResultTypes;
      status: StatusTypes;
      triggeredID: string;
    }
  ];
  state: 'triggered' | 'waiting' | 'suspended' | 'paused' | 'finished' | 'cancelled' | 'timedOut';
  stateBeforePause: string;
}
