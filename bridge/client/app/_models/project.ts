import semver from 'semver';
import { Stage } from './stage';
import { DeploymentInformation, Service } from './service';
import { IGitDataExtended, IProject } from '../../../shared/interfaces/project';

export function getShipyardVersion(project?: IProject): string {
  return project?.shipyardVersion?.split('/').pop() ?? '';
}

export function isShipyardNotSupported(
  project: IProject | undefined,
  supportedVersion: string | undefined | null
): boolean {
  const version = getShipyardVersion(project);
  return supportedVersion !== null && (!version || !supportedVersion || semver.lt(version, supportedVersion));
}

export function getDistinctServiceNames(project?: IProject): string[] {
  return project
    ? Array.from(
        project.stages.reduce((set, stage) => {
          stage.services.forEach((s) => {
            set.add(s.serviceName);
            return null;
          });
          return set;
        }, new Set<string>())
      )
    : [];
}

export class Project implements IProject {
  public projectName!: string;
  public gitCredentials?: IGitDataExtended;
  public shipyardVersion?: string;
  public projectDetailsLoaded = false; // true if project was fetched via project endpoint of bridge server
  public stages: Stage[] = [];
  public services?: Service[];
  public creationDate!: string;
  public shipyard!: string;

  static fromJSON(data: IProject): Project {
    const project: Project = Object.assign(new this(), data);
    project.stages = project.stages.map((stage) => Stage.fromJSON(stage));
    return project;
  }

  // returns a project without default values
  get reduced(): Partial<Project> {
    const { projectDetailsLoaded, ...copyProject } = this;
    return copyProject;
  }

  // replace project with a new one, but keep references
  public update(project: Project): void {
    this.gitCredentials = project.gitCredentials;
    const services: { [name: string]: Service } = {};
    for (const newStage of project.stages) {
      const existingStage = this.stages.find((stage) => stage.stageName === newStage.stageName);
      if (existingStage) {
        existingStage.update(newStage);
      }
      // at the moment deleting/adding stages is not supported, so we don't need to consider this case for now

      for (const service of newStage.services) {
        if (!services[service.serviceName]) {
          services[service.serviceName] = service;
        }
      }
      this.services = Object.values(services);
    }
  }

  getServices(stageName?: string): Service[] {
    if (!stageName) {
      if (!this.services) {
        let services: Service[] = [];
        for (const currentStage of this.stages) {
          services = services.concat(
            // eslint-disable-next-line @typescript-eslint/no-loop-func
            currentStage.services.filter((s: Service) => !services.some((ss) => ss.serviceName === s.serviceName))
          );
        }
        this.services = services;
      }
      return this.services;
    } else {
      return this.stages.find((s) => s.stageName === stageName)?.services ?? [];
    }
  }

  getServiceNames(): string[] {
    return this.services?.map((service) => service.serviceName) ?? [];
  }

  isShipyardNotSupported(supportedVersion: string | undefined | null): boolean {
    return isShipyardNotSupported(this, supportedVersion);
  }

  getService(serviceName: string): Service | undefined {
    return this.getServices().find((s) => s.serviceName === serviceName);
  }

  getLatestDeployment(service?: Service, stage?: Stage): DeploymentInformation | undefined {
    let currentService: Service | undefined;
    if (service) {
      if (stage) {
        currentService = this.getServices(stage.stageName)?.find((s) => s.serviceName === service.serviceName);
      } else {
        currentService = this.getService(service.serviceName);
      }
    }
    return currentService?.deploymentInformation;
  }

  public getStageNames(): string[] {
    return this.stages.map((stage) => stage.stageName);
  }

  public getStages(parent: string[] | null): Stage[] {
    return this.stages.filter(
      (s) => (parent && s.parentStages?.every((element, i) => element === parent[i])) || (!parent && !s.parentStages)
    );
  }

  public getParentStages(): [null, ...string[][]] {
    return this.stages.reduce(
      (stages: [null, ...string[][]], stage) => {
        if (
          stage.parentStages &&
          !stages.find((parent) => parent?.every((element, i) => element === stage.parentStages?.[i]))
        ) {
          stages.push(stage.parentStages);
        }
        return stages;
      },
      [null]
    );
  }
}
