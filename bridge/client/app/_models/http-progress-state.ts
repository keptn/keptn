export enum HttpProgressState {
  start,
  end
}

export interface HttpState {
  url: string;
  state: HttpProgressState;
}
