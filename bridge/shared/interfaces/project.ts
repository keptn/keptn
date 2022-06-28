import { IStage } from './stage';

export interface IProject {
  projectName: string;
  gitUser?: string;
  gitRemoteURI?: string;
  shipyardVersion?: string;
  gitProxyScheme?: 'https' | 'http';
  gitProxyUrl?: string;
  gitProxyUser?: string;
  gitProxyInsecure: boolean;
  stages: IStage[];
}
