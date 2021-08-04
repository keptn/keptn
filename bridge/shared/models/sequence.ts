import { EvaluationResult } from '../interfaces/evaluation-result';
import { EventState } from './event-state';

export type SequenceEvent = {
  id: string,
  time: string,
  type: string
};

export type SequenceStage = {
  image?: string,
  latestEvaluation?: EvaluationResult,
  latestEvent?: SequenceEvent,
  latestFailedEvent?: SequenceEvent,
  name: string,
};

enum sequenceState {
  WAITING = 'waiting'
}
export const SequenceState = {...sequenceState, ...EventState};
export type SequenceStateType = sequenceState & EventState;

export class Sequence {
  name!: string;
  project!: string;
  service!: string;
  shkeptncontext!: string;
  stages!: SequenceStage[];
  state!: SequenceStateType;
  time!: string;
  problemTitle?: string;
}
