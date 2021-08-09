import { Sequence } from './sequence';
import { SequenceResult as sr } from '../../../shared/interfaces/sequence-result';

export interface SequenceResult extends sr {
  states: Sequence[];
}
