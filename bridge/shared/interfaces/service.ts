export interface IServiceEvent {
  eventId: string;
  keptnContext: string;
  time: string; // nanoseconds
}

export interface IService {
  serviceName: string;
  creationDate: string; // nanoseconds
  deployedImage?: string;
  lastEventTypes?: { [event: string]: IServiceEvent | undefined };
}
