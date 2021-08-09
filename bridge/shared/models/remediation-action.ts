import { ResultTypes } from './result-types';

export class RemediationAction {
  action!: string;
  description!: string;
  name!: string;
  state!: 'triggered' | 'started' | 'finished';
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
