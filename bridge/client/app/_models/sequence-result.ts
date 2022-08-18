import { SequenceState } from './sequenceState';
import { SequenceResult as sr } from '../../../shared/interfaces/sequence-result';

export interface SequenceResult extends sr {
  states: SequenceState[];
}
