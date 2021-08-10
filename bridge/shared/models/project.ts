import { Stage } from './stage';
import { Service } from './service';

export class Project {
  projectName!: string;
  gitUser?: string;
  gitRemoteURI?: string;
  shipyardVersion?: string;
  stages: Stage[] = [];
  services?: Service[];
}
