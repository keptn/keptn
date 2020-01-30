export enum HttpProgressState {
  start,
  end
}

export class HttpState {
  url: string;
  state: HttpProgressState;
}
