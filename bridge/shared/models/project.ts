import { Stage } from './stage';

export class Project {
  projectName!: string;
  gitUser?: string;
  gitRemoteURI?: string;
  shipyardVersion?: string;
  stages: Stage[] = [];
}
