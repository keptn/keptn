import { ResultTypes } from './result-types';
import { EventState } from './event-state';

export type IRemediationAction = {
  action: string;
  description: string;
  name: string;
  state: EventState,
  result?: ResultTypes
};

export class RemediationAction implements IRemediationAction {
  action!: string;
  description!: string;
  name!: string;
  state!: EventState;
  result?: ResultTypes;

  public static fromJSON(data: unknown): RemediationAction {
    return Object.assign(new this(), data);
  }

  public getDetails(): string | undefined {
    return this.description || this.name;
  }

  public isFinished(): boolean {
    return this.state === 'finished';
  }
}
