import { Resource } from './resource';

export interface ServiceResource extends Resource {
  stageName: string;
}
