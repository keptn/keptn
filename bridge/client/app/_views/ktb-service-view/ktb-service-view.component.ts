import { Component, HostBinding, OnDestroy, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { combineLatest, Subject } from 'rxjs';
import { filter, map, switchMap, take, takeUntil } from 'rxjs/operators';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';
import { Location } from '@angular/common';
import { DeploymentInformation, ServiceState } from '../../../../shared/models/service-state';
import { AppUtils } from '../../_utils/app.utils';
import { DeploymentInformationSelection, DeploymentSelection } from '../../_interfaces/deployment-selection';

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
  public selectedDeploymentInfo?: DeploymentInformationSelection;
  public isQualityGatesOnly = false;
  public serviceStates?: ServiceState[];
  public deploymentInterval = 30_000;
  public projectName?: string;
  public selectedDeployment?: DeploymentSelection;
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
          const paramDeployment = serviceStates
            .find((state) => state.name === serviceName)
            ?.deployments.find((deployment) => deployment.keptnContext === keptnContext);
          // const changedDeployments =
          //   (this.selectedDeploymentInfo &&
          //     service?.deployments.filter(
          //       (deployment) => deployment.name === this.selectedDeploymentInfo?.deployment.name
          //     )) ??
          //   []; // the context of a deployment may change
          this.setDeploymentInfo(
            project.projectName,
            serviceName,
            this.getSelectedDeployment(project.projectName, serviceName, paramDeployment),
            params.get('stage') ?? undefined
          );
        }
      });

    project$.subscribe((project) => {
      this.projectName = project.projectName;
    });
  }

  private updateServiceStates(serviceStates: ServiceState[]): void {
    if (!this.serviceStates) {
      this.serviceStates = serviceStates;
    } else {
      for (const serviceState of serviceStates) {
        // deployments.length === 0 means that there aren't any updates for a service
        if (serviceState.deployments.length) {
          const serviceStateOriginal = this.serviceStates.find(
            (serviceStateO) => serviceStateO.name === serviceState.name
          );
          if (serviceStateOriginal) {
            for (const deploymentNew of serviceState.deployments) {
              const deploymentOriginal = serviceStateOriginal.deployments.find(
                (deployment) => deployment.keptnContext === deploymentNew.keptnContext
              );
              if (deploymentOriginal) {
                // update existing deployment
                deploymentOriginal.stages = [...deploymentOriginal.stages, ...deploymentNew.stages];

                // update other deployments (remove the stages)
                for (let i = 0; i < serviceStateOriginal.deployments.length; ++i) {
                  const deployment = serviceStateOriginal.deployments[i];
                  if (deployment !== deploymentOriginal) {
                    deployment.stages = deployment.stages.filter((stage) =>
                      deploymentNew.stages.some((st) => st.name === stage.name)
                    );
                    // delete deployment if it does not exist anymore
                    if (deployment.stages.length === 0) {
                      serviceStateOriginal.deployments.splice(i, 1);
                    }
                  }
                }
              } else {
                // add new deployment
                serviceStateOriginal.deployments.push(deploymentNew);
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

  private getSelectedDeployment(
    projectName: string,
    serviceName: string,
    paramDeployment?: DeploymentInformation
  ): DeploymentInformation | undefined {
    let selectedDeployment;
    if (paramDeployment) {
      selectedDeployment = paramDeployment;
    } else {
      const routeUrl = this.router.createUrlTree(['/project', projectName, 'service', serviceName]);
      this.location.go(routeUrl.toString());
    }
    return selectedDeployment;
  }

  private setDeploymentInfo(
    projectName: string,
    serviceName: string,
    selectedDeployment?: DeploymentInformation,
    paramStage?: string
  ): void {
    if (selectedDeployment) {
      let stage;
      if (paramStage) {
        stage = paramStage;
      } else {
        stage = selectedDeployment.stages[selectedDeployment.stages.length - 1].name;
        const routeUrl = this.router.createUrlTree([
          '/project',
          projectName,
          'service',
          serviceName,
          'context',
          selectedDeployment.keptnContext,
          'stage',
          stage,
        ]);
        this.location.go(routeUrl.toString());
      }
      this.selectedDeploymentInfo = { deployment: selectedDeployment, stage };
    } else {
      this.selectedDeploymentInfo = undefined;
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
        for (const deployment of serviceState.deployments) {
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
    return serviceState.deployments.some((deployment) => deployment.stages.some((stage) => stage.hasOpenRemediations));
  }

  public getLatestImage(serviceState: ServiceState): string | undefined {
    let latestTime: Date | undefined;
    let image: string | undefined;
    for (const deployment of serviceState.deployments) {
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
    this.deploymentLoading = true;
    this.dataService.getServiceDeployment(projectName, deploymentInfo.deployment.keptnContext).subscribe(
      (deployment) => {
        this.selectedDeployment = {
          deployment,
          selectedImage: deploymentInfo.deployment.image,
          stage: deploymentInfo.stage,
        };
        this.deploymentLoading = false;
      },
      () => {
        this.deploymentLoading = false;
      }
    );
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
