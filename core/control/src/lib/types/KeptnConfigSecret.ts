import { KeptnConfig } from './KeptnConfig';
export interface KeptnConfigSecret {
  apiVersion: string;
  kind: string;
  metadata: Metadata;
  type: string;
  data: KeptnConfig;
}
interface Metadata {
  name: string;
  namespace: string;
}
