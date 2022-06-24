import { IProject } from '../models/IProject';

export interface IProjectResult {
  projects: IProject[];
  totalCount: number;
}
