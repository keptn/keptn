import { Stage } from './stage';
import { Project as pj } from '../../shared/models/project';

export class Project extends pj {
  stages: Stage[] = [];

  public static fromJSON(data: unknown): Project {
    const project: Project = Object.assign(new this(), data);
    project.stages = project.stages.map(stage => Stage.fromJSON(stage));
    return project;
  }
}
