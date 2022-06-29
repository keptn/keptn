export enum Features {
  ROOT = 'root',
}

export const enum LoadingState {
  INIT = 'INIT',
  LOADING = 'LOADING',
  LOADED = 'LOADED',
}
export interface ErrorState {
  errorMsg: string;
}

export type CallState = LoadingState | ErrorState;

export interface ApiCall<T> {
  data: T;
  call: CallState;
}
