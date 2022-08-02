export enum DeleteType {
  PROJECT = 'project',
  SERVICE = 'service',
  SUBSCRIPTION = 'subscription',
}

export enum DeleteResult {
  ERROR = 'error',
  SUCCESS = 'success',
}

export interface DeleteData {
  type: DeleteType;
  name?: string;
  context?: unknown;
}

export interface DeletionTriggeredEvent {
  type: DeleteType;
  name?: string;
  context?: unknown;
}

export interface DeletionProgressEvent {
  result?: DeleteResult;
  isInProgress: boolean;
  error?: string;
}
