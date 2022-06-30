import { getDistinctServiceNames, Project } from './project';
import { waitForAsync } from '@angular/core/testing';
import { IProject } from '../../../shared/interfaces/project';

const projects: IProject[] = [
  {
    projectName: 'sockshop',
    creationDate: '',
    shipyard: '',
    stages: [
      {
        services: [
          { serviceName: 'carts', creationDate: '' },
          { serviceName: 'carts-db', creationDate: '' },
        ],
        stageName: 'dev',
      },
      {
        services: [
          { serviceName: 'carts', creationDate: '' },
          { serviceName: 'carts-db', creationDate: '' },
        ],
        stageName: 'staging',
      },
      {
        services: [
          { serviceName: 'carts', creationDate: '' },
          { serviceName: 'carts-db', creationDate: '' },
          { serviceName: 'carts-db2', creationDate: '' },
        ],
        stageName: 'production',
      },
    ],
  },
];

describe('Project', () => {
  it('should create instances from json', waitForAsync(() => {
    const projectsClasses: Project[] = projects.map((project) => Project.fromJSON(project));

    expect(projectsClasses[0]).toBeInstanceOf(Project);
    expect(projectsClasses[0].projectName).toEqual('sockshop');
    expect(projectsClasses[0].getServices().length).toEqual(3);
    expect(projectsClasses[0].stages.length).toEqual(3);
  }));

  it('should be a supported shipyard version if provided version is null', () => {
    // if not provided (null) we assume it's supported (metadata endpoint not available)
    const project = createProjectWithShipyard('spec.keptn.sh/0.2.0');
    expect(project.isShipyardNotSupported(null)).toBe(false);
  });

  it('should be a supported shipyard version if provided version greater', () => {
    const project = createProjectWithShipyard('spec.keptn.sh/0.2.0');
    expect(project.isShipyardNotSupported('0.1.1')).toBe(false);
  });

  it('should be a not supported shipyard version if provided version is lower', () => {
    const project = createProjectWithShipyard('spec.keptn.sh/0.2.0');
    expect(project.isShipyardNotSupported('0.3.0')).toBe(true);
  });

  it('should be a not supported shipyard version if provided version is empty', () => {
    const project = createProjectWithShipyard('spec.keptn.sh/0.2.0');
    expect(project.isShipyardNotSupported('')).toBe(true);
  });

  it('should be a not supported shipyard version if provided version is undefined', () => {
    const project = createProjectWithShipyard('spec.keptn.sh/0.2.0');
    expect(project.isShipyardNotSupported(undefined)).toBe(true);
  });

  it('should be a not supported shipyard if project does not have a shipyard version', () => {
    const project = new Project();
    expect(project.isShipyardNotSupported('10.10.10')).toBe(true);
    expect(project.isShipyardNotSupported('0.0.0')).toBe(true);
  });

  it('should be a supported shipyard version if the provided version is the same', () => {
    const project = createProjectWithShipyard('spec.keptn.sh/0.2.0');
    expect(project.isShipyardNotSupported('0.2.0')).toBe(false);
  });

  describe(getDistinctServiceNames.name, () => {
    it('should return distinct services (deprecated)', () => {
      // given
      const projectsClasses: Project[] = projects.map((project) => Project.fromJSON(project));
      projectsClasses.forEach((p) => p.getServices());

      // when
      const serviceNames = projectsClasses[0].getServiceNames();

      // then
      expect(serviceNames).toEqual(['carts', 'carts-db', 'carts-db2']);
    });

    it('should return distinct services', () => {
      // given
      const project = projects[0];

      // when
      const serviceNames = getDistinctServiceNames(project);

      // then
      expect(serviceNames).toEqual(['carts', 'carts-db', 'carts-db2']);
    });
  });

  function createProjectWithShipyard(shipyardVersion: string): Project {
    const project = new Project();
    project.shipyardVersion = shipyardVersion;
    return project;
  }
});
