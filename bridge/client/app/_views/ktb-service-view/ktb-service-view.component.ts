import { Component, HostBinding, Inject, OnDestroy, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { combineLatest, Subject } from 'rxjs';
import { filter, map, skip, switchMap, take, takeUntil } from 'rxjs/operators';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';
import { Location } from '@angular/common';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { DeploymentInformationSelection } from '../../_interfaces/deployment-selection';
import { DeploymentInformation, ServiceState } from '../../_models/service-state';
import { SequenceState } from '../../../../shared/models/sequence';

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
  public projectName?: string;
  public selectedDeployment?: DeploymentInformationSelection;
  public deploymentLoading = false;

  constructor(
    private dataService: DataService,
    private route: ActivatedRoute,
    private router: Router,
    public location: Location,
    @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number
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

    params$.pipe(take(1)).subscribe((params) => {
      this.serviceName = params.get('serviceName') ?? undefined;
    });

    const serviceStates$ = projectName$.pipe(
      switchMap((projectName) => AppUtils.createTimer(0, this.initialDelayMillis).pipe(map(() => projectName))),
      takeUntil(this.unsubscribe$),
      switchMap((projectName) =>
        this.dataService.getServiceStates(projectName, this.getLatestDeploymentTime()?.toISOString())
      ),
      takeUntil(this.unsubscribe$)
    );

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
    serviceStates$.pipe(skip(1)).subscribe((serviceStates) => {
      this.updateServiceStates(serviceStates);
    });

    project$.subscribe((project) => {
      this.projectName = project.projectName;
    });
  }

  // checks if the given stage exists in the deployment and returns the latest one if not
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
      ServiceState.update(this.serviceStates, serviceStates);
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

  public getLatestImage(serviceState: ServiceState): string {
    let latestTime: Date | undefined;
    let image = 'unknown';
    for (const deployment of serviceState.deploymentInformation) {
      const latestStageTime = deployment.stages.reduce((max: undefined | Date, stage) => {
        const date = new Date(stage.time);
        return max && max > date ? max : date;
      }, undefined);
      if (deployment.image && latestStageTime && (!latestTime || latestStageTime > latestTime)) {
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
      const originalDeployment = deploymentInfo.deploymentInformation.deployment;
      if (
        this.projectName &&
        (originalDeployment.state === SequenceState.FINISHED || originalDeployment.state === SequenceState.TIMEDOUT)
      ) {
        this.dataService
          .getOpenRemediationsOfService(this.projectName, originalDeployment.service)
          .subscribe((remediations) => {
            originalDeployment.updateRemediations(remediations);
          });
      } else {
        this.dataService
          .getServiceDeployment(
            projectName,
            deploymentInfo.deploymentInformation.keptnContext,
            new Date(originalDeployment.latestTimeUpdated).toISOString()
          )
          .subscribe(
            (deployment) => {
              originalDeployment.update(deployment);
              this.selectedDeployment = deploymentInfo;
              this.deploymentLoading = false;
            },
            () => {
              this.deploymentLoading = false;
            }
          );
      }
    }
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
