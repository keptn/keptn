import { Stage } from './stage';

export class Project {
  projectName!: string;
  gitUser?: string;
  gitRemoteURI?: string;
  shipyardVersion?: string;
  gitProxyScheme?: 'https' | 'http';
  gitProxyUrl?: string;
  gitProxyUser?: string;
  gitProxyInsecure = false;
  stages: Stage[] = [];
}
