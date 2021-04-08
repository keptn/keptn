import {Project} from './project';
import {Stage} from './stage';
import {Service} from './service';

export const ProjectMock = {
  projectName: 'sockshop',
  stages: [
    {
      stageName: 'development'
    } as Stage,
    {
      stageName: 'staging'
    } as Stage,
    {
      stageName: 'production'
    } as Stage
  ],
  services: [
    {
      serviceName: 'cards'
    } as Service,
    {
      serviceName: 'cards-db'
    } as Service
  ]
} as Project;
