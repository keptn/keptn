export interface IServiceEvent {
  eventId: string;
  keptnContext: string;
  time: string; // nanoseconds
}

export interface IService {
  serviceName: string;
  creationDate: number;
  stage: string;
  deployedImage?: string;
  lastEventTypes?: { [p: string]: IServiceEvent | undefined };
}
