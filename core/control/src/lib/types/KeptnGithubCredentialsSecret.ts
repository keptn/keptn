import { KeptnGithubCredentials } from './KeptnGithubCredentials';

export interface KeptnGithubCredentialsSecret {
  apiVersion: string;
  kind: string;
  metadata: Metadata;
  type: string;
  data: KeptnGithubCredentials;
}

interface Metadata {
  name: string;
  namespace: string;
}
