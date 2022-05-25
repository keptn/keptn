import { Project } from './project';
import { waitForAsync } from '@angular/core/testing';

describe('Project', () => {
  it('should create instances from json', waitForAsync(() => {
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

  function createProjectWithShipyard(shipyardVersion: string): Project {
    const project = new Project();
    project.shipyardVersion = shipyardVersion;
    return project;
  }
});
