import { Trace as tc } from '../../shared/models/trace';

export class Trace extends tc {
  static fromJSON(data: unknown): Trace {
    return Object.assign(new this(), data);
  }
}
