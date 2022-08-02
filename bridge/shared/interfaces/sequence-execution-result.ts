import { IShipyardSequence } from './shipyard';

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
  labels: {
    additionalProp1: string;
    additionalProp2: string;
    additionalProp3: string;
  };
  message: string;
  project: string;
  result: string;
  service: string;
  stage: string;
  status: string;
  triggeredId: string;
}

export interface SequenceExecutionStatus {
  currentTask: {
    events: [
      {
        eventType: string;
        properties: {
          // eslint-disable-next-line @typescript-eslint/ban-types
          additionalProp1: {};
        };
        result: string;
        source: string;
        status: string;
        time: string;
      }
    ];
    name: string;
    triggeredID: string;
  };
  previousTasks: [
    {
      name: string;
      properties: {
        // eslint-disable-next-line @typescript-eslint/ban-types
        additionalProp1: {};
      };
      result: string;
      status: string;
      triggeredID: string;
    }
  ];
  state: string;
  stateBeforePause: string;
}
