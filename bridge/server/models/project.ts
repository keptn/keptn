import { Stage } from './stage';
import { IGitDataExtended, IProject } from '../../shared/interfaces/Project';

export class Project implements IProject {
  public gitCredentials?: IGitDataExtended;
  public projectName = '';
  public shipyardVersion?: string;
  public stages: Stage[] = [];
  public creationDate = '';
  public shipyard = '';

  public static fromJSON(data: unknown): Project {
    const project: Project = Object.assign(new this(), data);
    project.stages = project.stages.map((stage) => Stage.fromJSON(stage));
    return project;
  }
}
