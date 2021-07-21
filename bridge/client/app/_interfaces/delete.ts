export enum DeleteType {
  PROJECT = 'project'
}

export enum DeleteResult {
  ERROR = 'error',
  SUCCESS = 'success'
}

export interface DeleteData {
  type: DeleteType;
  name: string;
}

export interface DeletionTriggeredEvent {
  type: DeleteType;
  name: string;
}

export interface DeletionProgressEvent {
  result: DeleteResult | null;
  isInProgress: boolean;
  error: string | null;
}
