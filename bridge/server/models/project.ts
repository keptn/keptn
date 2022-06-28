import { Stage } from './stage';
import { IProject } from '../../shared/interfaces/project';

export class Project implements IProject {
  stages: Stage[] = [];

  public static fromJSON(data: unknown): Project {
    const project: Project = Object.assign(new this(), data);
    project.stages = project.stages.map((stage) => Stage.fromJSON(stage));
    return project;
  }

  gitProxyInsecure = false;
  projectName = '';
}
