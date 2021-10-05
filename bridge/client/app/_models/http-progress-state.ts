export enum HttpProgressState {
  START,
  END,
}

export interface HttpState {
  url: string;
  state: HttpProgressState;
}
