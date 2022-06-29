import { SequenceResult as sr } from '../../../shared/interfaces/sequence-result';
import { ISequence } from '../../../shared/interfaces/sequence';

export interface SequenceResult extends sr {
  states: ISequence[];
}
