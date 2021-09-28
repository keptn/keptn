import { ProjectResult as pr } from '../../../shared/interfaces/project-result';
import { Project } from '../_models/project';

export interface ProjectResult extends pr {
  projects: Project[];
  totalCount: number;
}
