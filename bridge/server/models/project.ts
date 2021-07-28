import { ProjectResult } from '../interfaces/project-result.js';
import { Stage } from './stage.js';

export class Project implements ProjectResult {
  creationDate!: number;
  projectName!: string;
  shipyard?: string;
  shipyardVersion?: string;
  stages: Stage[] = [];

  public static fromJSON(data: unknown): Project {
    const project: Project = Object.assign(new this(), data);
    project.stages = project.stages.map(stage => Stage.fromJSON(stage));
    return project;
  }
}
