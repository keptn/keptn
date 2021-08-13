import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { combineLatest, Subject } from 'rxjs';
import { filter, switchMap, take, takeUntil } from 'rxjs/operators';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';
import { Deployment } from 'client/app/_models/deployment';
import { Location } from '@angular/common';

@Component({
  selector: 'ktb-service-view',
  templateUrl: './ktb-service-view.component.html',
  styleUrls: ['./ktb-service-view.component.scss'],
  host: {
    class: 'ktb-service-view'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbServiceViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public project?: Project;
  public serviceName?: string;
  public selectedDeploymentInfo?: { deployment: Deployment, stage: string};
  public isQualityGatesOnly = false;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute, private router: Router, private location: Location) { }

  ngOnInit() {
    this.dataService.isQualityGatesOnly.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(isQualityGatesOnly => {
      this.isQualityGatesOnly = isQualityGatesOnly;
    });

    this.dataService.changedDeployments
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this._changeDetectorRef.markForCheck();
      });

    const params$ = this.route.params
      .pipe(takeUntil(this.unsubscribe$));

    const project$ = params$.pipe(
      switchMap(params => this.dataService.getProject(params.projectName)),
      filter((project: Project | undefined): project is Project => !!project),
      takeUntil(this.unsubscribe$)
      );

    params$.pipe(take(1)).subscribe(params => {
      this.serviceName = params.serviceName;
      this._changeDetectorRef.markForCheck();
    });

    combineLatest([params$, project$]).pipe(take(1)).subscribe(([params, project]) => {
      if (params.shkeptncontext && params.serviceName) {
        const service = project.getServices().find(s => s.serviceName === params.serviceName);
        const paramDeployment = service?.deployments.find(deployment => deployment.shkeptncontext === params.shkeptncontext);
        const changedDeployments = (this.selectedDeploymentInfo && service?.deployments.filter(deployment => deployment.name === this.selectedDeploymentInfo?.deployment.name)) ?? []; // the context of a deployment may change
        this.setDeploymentInfo(project.projectName, this.getSelectedDeployment(project.projectName, params.serviceName, changedDeployments, paramDeployment), params.stage);
      }
    });

    project$.subscribe(project => {
      if (this.selectedDeploymentInfo) { // the selected deployment gets lost if the project is updated, because the deployments are rebuild
        const selectedDeployment = project.getServices().find(s => s.serviceName === this.selectedDeploymentInfo?.deployment.service)?.deployments.find(d => d.shkeptncontext === this.selectedDeploymentInfo?.deployment.shkeptncontext);
        if (selectedDeployment) {
          this.setDeploymentInfo(project.projectName, selectedDeployment, this.selectedDeploymentInfo.stage);
        }
      }
      this.dataService.loadOpenRemediations(project);
      this.project = project;
      this._changeDetectorRef.markForCheck();
    });
  }

  private getSelectedDeployment(projectName: string, serviceName: string, changedDeployments: Deployment[], paramDeployment?: Deployment): Deployment | undefined {
    let selectedDeployment;
    if (paramDeployment) {
      selectedDeployment = paramDeployment;
    } else if (changedDeployments.length > 0) {
      if (changedDeployments.length === 1) {
        selectedDeployment = changedDeployments[0];
      } else {
        selectedDeployment = changedDeployments.find(d => d.stages.some(s => this.selectedDeploymentInfo?.deployment.stages.some(sd => s.stageName === sd.stageName)));
      }
    } else {
      const routeUrl = this.router.createUrlTree(['/project', projectName, 'service', serviceName]);
      this.location.go(routeUrl.toString());
    }
    return selectedDeployment;
  }

  private setDeploymentInfo(projectName: string, selectedDeployment?: Deployment, paramStage?: string) {
    if (selectedDeployment) {
      let stage;
      if (paramStage && selectedDeployment.hasStage(paramStage)) {
        stage = paramStage;
      }
      else {
        stage = selectedDeployment.stages[selectedDeployment.stages.length - 1].stageName;
        const routeUrl = this.router.createUrlTree(['/project', projectName, 'service', selectedDeployment.service, 'context', selectedDeployment.shkeptncontext, 'stage', stage]);
        this.location.go(routeUrl.toString());
      }
      this.selectedDeploymentInfo = {deployment: selectedDeployment, stage};
    }
    else {
      this.selectedDeploymentInfo = undefined;
    }
  }

  public selectService(projectName: string, serviceName: string): void {
    if (this.serviceName !== serviceName) {
      this.serviceName = serviceName;
      this._changeDetectorRef.markForCheck();
    }
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
