import { ProjectResult as pr } from '../../shared/interfaces/project-result';
import { Project } from '../models/project';

export interface ProjectResult extends pr {
  projects: Project[];
}
