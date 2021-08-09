import { EventResult as er } from '../../shared/interfaces/event-result';
import { Trace } from '../models/trace';

export interface EventResult extends er {
  events: Trace[];
}
