import { Component, HostBinding, OnDestroy, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { combineLatest, Subject } from 'rxjs';
import { filter, map, switchMap, take, takeUntil } from 'rxjs/operators';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';
import { Location } from '@angular/common';
import { ServiceState } from '../../../../shared/models/service-state';
import { AppUtils } from '../../_utils/app.utils';
import { DeploymentInformationSelection } from '../../_interfaces/deployment-selection';
import { DeploymentInformation } from '../../_models/service-state';

@Component({
  selector: 'ktb-service-view',
  templateUrl: './ktb-service-view.component.html',
  styleUrls: ['./ktb-service-view.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
})
export class KtbServiceViewComponent implements OnDestroy {
  @HostBinding('class') cls = 'ktb-service-view';
  private readonly unsubscribe$ = new Subject<void>();
  public serviceName?: string;
  public isQualityGatesOnly = false;
  public serviceStates?: ServiceState[];
  public deploymentInterval = 30_000;
  public projectName?: string;
  public selectedDeployment?: DeploymentInformationSelection;
  public deploymentLoading = false;

  constructor(
    private dataService: DataService,
    private route: ActivatedRoute,
    private router: Router,
    public location: Location
  ) {
    this.dataService.isQualityGatesOnly.pipe(takeUntil(this.unsubscribe$)).subscribe((isQualityGatesOnly) => {
      this.isQualityGatesOnly = isQualityGatesOnly;
    });

    const params$ = this.route.paramMap.pipe(takeUntil(this.unsubscribe$));
    const projectName$ = params$.pipe(
      map((params) => params.get('projectName')),
      filter((projectName: string | null): projectName is string => !!projectName)
    );

    const project$ = projectName$.pipe(
      switchMap((projectName) => this.dataService.getProject(projectName)),
      filter((project: Project | undefined): project is Project => !!project?.projectDetailsLoaded),
      takeUntil(this.unsubscribe$)
    );

    const serviceStates$ = projectName$.pipe(
      switchMap((projectName) => this.dataService.getServiceStates(projectName))
    );

    params$.pipe(take(1)).subscribe((params) => {
      this.serviceName = params.get('serviceName') ?? undefined;
    });

    projectName$
      .pipe(
        switchMap((projectName) =>
          AppUtils.createTimer(this.deploymentInterval, this.deploymentInterval).pipe(map(() => projectName))
        ),
        takeUntil(this.unsubscribe$),
        switchMap((projectName) =>
          this.dataService.getServiceStates(projectName, this.getLatestDeploymentTime()?.toISOString())
        ),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((serviceStates) => {
        this.updateServiceStates(serviceStates);
      });

    combineLatest([params$, project$, serviceStates$])
      .pipe(take(1))
      .subscribe(([params, project, serviceStates]) => {
        this.updateServiceStates(serviceStates);
        const keptnContext = params.get('shkeptncontext');
        const serviceName = params.get('serviceName');
        if (keptnContext && serviceName) {
          const selectedDeploymentInformation = serviceStates
            .find((state) => state.name === serviceName)
            ?.deploymentInformation.find((deployment) => deployment.keptnContext === keptnContext);
          if (selectedDeploymentInformation) {
            const stage = this.validateStage(
              selectedDeploymentInformation,
              project.projectName,
              serviceName,
              params.get('stage')
            );

            const selection = {
              deploymentInformation: selectedDeploymentInformation,
              stage,
            };
            this.deploymentSelected(selection, project.projectName);
          }
        }
      });

    project$.subscribe((project) => {
      this.projectName = project.projectName;
    });
  }

  private validateStage(
    selectedDeploymentInformation: DeploymentInformation,
    projectName: string,
    serviceName: string,
    stage: string | null
  ): string {
    if (!stage || selectedDeploymentInformation.stages.some((s) => s.name === stage)) {
      stage = selectedDeploymentInformation.stages[selectedDeploymentInformation.stages.length - 1].name;
      const routeUrl = this.router.createUrlTree([
        '/project',
        projectName,
        'service',
        serviceName,
        'context',
        selectedDeploymentInformation.keptnContext,
        'stage',
        stage,
      ]);
      this.location.go(routeUrl.toString());
    }
    return stage;
  }

  private updateServiceStates(serviceStates: ServiceState[]): void {
    if (!this.serviceStates) {
      this.serviceStates = serviceStates;
    } else {
      for (const serviceState of serviceStates) {
        // deployments.length === 0 means that there aren't any updates for a service
        if (serviceState.deploymentInformation.length) {
          const serviceStateOriginal = this.serviceStates.find(
            (serviceStateO) => serviceStateO.name === serviceState.name
          );
          if (serviceStateOriginal) {
            for (const deploymentNew of serviceState.deploymentInformation) {
              const deploymentOriginal = serviceStateOriginal.deploymentInformation.find(
                (deployment) => deployment.keptnContext === deploymentNew.keptnContext
              );
              if (deploymentOriginal) {
                // update existing deployment
                deploymentOriginal.stages = [...deploymentOriginal.stages, ...deploymentNew.stages];

                // update other deployments (remove the stages)
                for (let i = 0; i < serviceStateOriginal.deploymentInformation.length; ++i) {
                  const deployment = serviceStateOriginal.deploymentInformation[i];
                  if (deployment !== deploymentOriginal) {
                    deployment.stages = deployment.stages.filter((stage) =>
                      deploymentNew.stages.some((st) => st.name === stage.name)
                    );
                    // delete deployment if it does not exist anymore
                    if (deployment.stages.length === 0) {
                      serviceStateOriginal.deploymentInformation.splice(i, 1);
                    }
                  }
                }
              } else {
                // add new deployment
                serviceStateOriginal.deploymentInformation.push(deploymentNew);
              }
            }
          } else {
            // new service with deployments
            this.serviceStates.push(serviceState);
          }
        } else if (!this.serviceStates.some((s) => s.name === serviceState.name)) {
          // new service
          this.serviceStates.push(serviceState);
        }
        // remove deleted services
        for (let i = 0; i < this.serviceStates.length; ++i) {
          const serviceStateOriginal = this.serviceStates[i];
          if (!serviceStates.some((state) => state.name === serviceStateOriginal.name)) {
            this.serviceStates.splice(i, 1);
          }
        }
      }
    }
  }

  public selectService(projectName: string, serviceName: string): void {
    if (this.serviceName !== serviceName) {
      this.serviceName = serviceName;
    }
  }

  private getLatestDeploymentTime(): Date | undefined {
    let latestTime: undefined | Date;
    if (this.serviceStates) {
      for (const serviceState of this.serviceStates) {
        for (const deployment of serviceState.deploymentInformation) {
          for (const stage of deployment.stages) {
            const date = new Date(stage.time);
            if (!latestTime || date < latestTime) {
              latestTime = date;
            }
          }
        }
      }
    }
    return latestTime;
  }

  public hasServiceRemediations(serviceState: ServiceState): boolean {
    return serviceState.deploymentInformation.some((deployment) =>
      deployment.stages.some((stage) => stage.hasOpenRemediations)
    );
  }

  public getLatestImage(serviceState: ServiceState): string | undefined {
    let latestTime: Date | undefined;
    let image: string | undefined;
    for (const deployment of serviceState.deploymentInformation) {
      const latestStageTime = deployment.stages.reduce((max: undefined | Date, stage) => {
        const date = new Date(stage.time);
        return max && max > date ? max : date;
      }, undefined);
      if (latestStageTime && (!latestTime || latestStageTime > latestTime)) {
        image = `${deployment.image}:${deployment.version}`;
        latestTime = latestStageTime;
      }
    }
    return image;
  }

  public deploymentSelected(deploymentInfo: DeploymentInformationSelection, projectName: string): void {
    if (!deploymentInfo.deploymentInformation.deployment) {
      this.deploymentLoading = true;
      this.dataService.getServiceDeployment(projectName, deploymentInfo.deploymentInformation.keptnContext).subscribe(
        (deployment) => {
          deploymentInfo.deploymentInformation.deployment = deployment;
          this.selectedDeployment = deploymentInfo;
          this.deploymentLoading = false;
        },
        () => {
          this.deploymentLoading = false;
        }
      );
    } else {
      // TODO: update
      this.selectedDeployment = deploymentInfo;
    }
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
