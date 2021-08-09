import { Project } from '../models/project';

export interface ProjectResult {
  projects: Project[];
  totalCount: number;
}
