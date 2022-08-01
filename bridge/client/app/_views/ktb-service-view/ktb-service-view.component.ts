import { Component, HostBinding, Inject, OnDestroy, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { combineLatest, Observable, Subject, Subscription } from 'rxjs';
import { filter, map, switchMap, takeUntil, tap } from 'rxjs/operators';
import { DataService } from '../../_services/data.service';
import { Location } from '@angular/common';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { DeploymentInformationSelection } from '../../_interfaces/deployment-selection';
import { ServiceState } from '../../_models/service-state';
import { Deployment } from '../../_models/deployment';
import { ServiceRemediationInformation } from '../../_models/service-remediation-information';

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
  public updateDeploymentSubscription$?: Subscription;
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

    const serviceStatesInitial$ = projectName$.pipe(
      tap(() => {
        // service states of another project are loaded; change to loading state
        this.serviceStates = undefined;
      }),
      switchMap((projectName) => this.dataService.getServiceStates(projectName)),
      takeUntil(this.unsubscribe$)
    );

    combineLatest([params$, projectName$, serviceStatesInitial$]).subscribe(([params, projectName, serviceStates]) => {
      this.updateServiceStates(serviceStates);
      const keptnContext = params.get('shkeptncontext');
      this.serviceName = params.get('serviceName') ?? undefined;
      if (keptnContext && this.serviceName) {
        const serviceState = serviceStates.find((state) => state.name === this.serviceName);
        const selectedDeploymentInformation = serviceState?.deploymentInformation.find(
          (deployment) => deployment.keptnContext === keptnContext
        );
        if (selectedDeploymentInformation) {
          const selection = {
            deploymentInformation: selectedDeploymentInformation,
            stage: params.get('stage') ?? '',
          };
          this.deploymentSelected(selection, projectName);
        } else if (serviceState) {
          // remove context and stage parameter if it does not exist
          const routeUrl = this.router.createUrlTree(['/project', projectName, 'service', serviceState.name]);
          this.location.go(routeUrl.toString());
        } else {
          // remove service parameter, if it does not exist
          const routeUrl = this.router.createUrlTree(['/project', projectName, 'service']);
          this.location.go(routeUrl.toString());
        }
      }
    });

    projectName$.subscribe((projectName) => {
      this.projectName = projectName;
      this.selectedDeployment = undefined;
    });

    if (this.initialDelayMillis !== 0) {
      const serviceStateInterval$ = projectName$.pipe(
        switchMap((projectName) =>
          AppUtils.createTimer(this.initialDelayMillis, this.initialDelayMillis).pipe(map(() => projectName))
        ),
        takeUntil(this.unsubscribe$),
        switchMap((projectName) => this.dataService.getServiceStates(projectName)),
        takeUntil(this.unsubscribe$)
      );
      serviceStateInterval$.subscribe((serviceStates) => {
        this.updateServiceStates(serviceStates);
      });
    }
  }

  // checks if the given stage exists in the deployment and returns the latest one if not
  private validateStage(deployment: Deployment, projectName: string, stage: string): string {
    if (!stage || !deployment.getStage(stage)) {
      stage = deployment.stages[deployment.stages.length - 1].name;
      const routeUrl = this.router.createUrlTree([
        '/project',
        projectName,
        'service',
        deployment.service,
        'context',
        deployment.keptnContext,
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

  public deploymentSelected(deploymentInfo: DeploymentInformationSelection, projectName: string): void {
    this.updateDeploymentSubscription$?.unsubscribe();
    this.updateDeploymentSubscription$ = AppUtils.createTimer(0, this.initialDelayMillis)
      .pipe(switchMap(() => this.updateDeployment(deploymentInfo, projectName)))
      .subscribe(
        (update) => {
          const originalDeployment = deploymentInfo.deploymentInformation.deployment;
          if (update instanceof Deployment) {
            if (!originalDeployment) {
              deploymentInfo.stage = this.validateStage(update, projectName, deploymentInfo.stage);
              deploymentInfo.deploymentInformation.deployment = update;
            } else {
              originalDeployment.update(update);
            }
            this.deploymentLoading = false;
            return;
          }
          if (originalDeployment && update) {
            originalDeployment.updateRemediations(update);
          }
        },
        () => {
          this.deploymentLoading = false;
        }
      );
  }

  private updateDeployment(
    deploymentInfo: DeploymentInformationSelection,
    projectName: string
  ): Observable<Deployment | ServiceRemediationInformation> {
    const originalDeployment = deploymentInfo.deploymentInformation.deployment;
    this.selectedDeployment = deploymentInfo;

    if (!originalDeployment) {
      // initially fetch deployment
      this.deploymentLoading = true;
      return this.dataService.getServiceDeployment(projectName, deploymentInfo.deploymentInformation.keptnContext);
    }
    // update deployment
    if (originalDeployment.isFinished()) {
      // deployment is finished. Just update open remediations
      return this.dataService.getOpenRemediationsOfService(projectName, originalDeployment.service);
    }

    return this.dataService.getServiceDeployment(
      projectName,
      deploymentInfo.deploymentInformation.keptnContext,
      originalDeployment.latestTimeUpdated?.toISOString()
    );
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
    this.updateDeploymentSubscription$?.unsubscribe();
  }
}
