import { Project } from './project';
import { waitForAsync } from '@angular/core/testing';

describe('Project', () => {
  it(
    'should create instances from json',
    waitForAsync(() => {
      const projects: Project[] = [
        {
          projectName: 'sockshop',
          stages: [
            { services: [{ serviceName: 'carts' }, { serviceName: 'carts-db' }], stageName: 'dev' },
            { services: [{ serviceName: 'carts' }, { serviceName: 'carts-db' }], stageName: 'staging' },
            { services: [{ serviceName: 'carts' }, { serviceName: 'carts-db' }], stageName: 'production' },
          ],
        },
      ].map((project) => Project.fromJSON(project));

      expect(projects[0]).toBeInstanceOf(Project);
      expect(projects[0].projectName).toEqual('sockshop');
      expect(projects[0].getServices().length).toEqual(2);
      expect(projects[0].stages.length).toEqual(3);
    })
  );
});
