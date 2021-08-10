import { SequenceResult as sr } from '../../shared/interfaces/sequence-result';
import { Sequence } from '../models/sequence';

export interface SequenceResult extends sr {
  states: Sequence[];
}
