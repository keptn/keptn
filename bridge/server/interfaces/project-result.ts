import { Stage } from '../models/stage.js';


export interface ProjectResult {
  creationDate: number;
  projectName: string;
  shipyard?: string;
  shipyardVersion?: string;
  stages: Stage[];
  gitRemoteUri?: string;
  gitUser?: string;
}
