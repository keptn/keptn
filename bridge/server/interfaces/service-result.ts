export interface ServiceResult {
  serviceName: string;
  creationDate: number;
  deployedImage?: string;
  lastEventTypes?: {[key: string]: {eventId: string, keptnContext: string, time: number}};
}
