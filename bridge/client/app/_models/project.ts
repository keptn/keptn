import semver from 'semver';
import { Stage } from './stage';
import { DeploymentInformation, Service } from './service';
import { IGitDataExtended } from '../_interfaces/git-upstream';
import { isGitInputWithHTTPS } from '../_utils/git-upstream.utils';
import { IProject } from '../../../shared/interfaces/project';

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
        project.stages.reduce((currentSet, stage) => {
          const serviceNames = stage.services.map((s) => s.serviceName);
          return new Set(...currentSet, ...serviceNames);
        }, new Set<string>())
      )
    : [];
}

export class Project implements IProject {
  projectName = '';
  gitUser?: string | undefined;
  gitRemoteURI?: string | undefined;
  shipyardVersion?: string | undefined;
  gitProxyScheme?: 'https' | 'http' | undefined;
  gitProxyUrl?: string | undefined;
  gitProxyUser?: string | undefined;
  gitProxyInsecure = false;
  private _gitUpstream?: IGitDataExtended;
  public projectDetailsLoaded = false; // true if project was fetched via project endpoint of bridge server
  public stages: Stage[] = [];
  public services?: Service[];

  static fromJSON(data: unknown): Project {
    const project: Project = Object.assign(new this(), data);
    project.stages = project.stages.map((stage) => Stage.fromJSON(stage));
    return project;
  }

  // returns a project without default values
  get reduced(): Partial<Project> {
    const { projectDetailsLoaded, ...copyProject } = this;
    return copyProject;
  }

  public get gitUpstream(): IGitDataExtended {
    if (!this._gitUpstream) {
      if (isGitInputWithHTTPS(this)) {
        this._gitUpstream = {
          https: {
            gitUser: this.gitUser,
            gitRemoteURL: this.gitRemoteURI ?? '',
            gitToken: '',
            gitProxyScheme: this.gitProxyScheme ?? 'https',
            gitProxyUrl: this.gitProxyUrl ?? '',
            gitProxyPassword: '',
            gitProxyUser: this.gitProxyUser ?? '',
            gitProxyInsecure: this.gitProxyInsecure,
          },
        };
      } else {
        this._gitUpstream = {
          ssh: {
            gitRemoteURL: this.gitRemoteURI ?? '',
            gitUser: this.gitUser,
            gitPrivateKey: '',
            gitPrivateKeyPass: '',
          },
        };
      }
    }
    return this._gitUpstream;
  }

  // replace project with a new one, but keep references
  public update(project: Project): void {
    this.gitRemoteURI = project.gitRemoteURI;
    this.gitUser = project.gitUser;
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

  getShipyardVersion(): string {
    return getShipyardVersion(this);
  }

  isShipyardNotSupported(supportedVersion: string | undefined | null): boolean {
    const version = this.getShipyardVersion();
    return supportedVersion !== null && (!version || !supportedVersion || semver.lt(version, supportedVersion));
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
