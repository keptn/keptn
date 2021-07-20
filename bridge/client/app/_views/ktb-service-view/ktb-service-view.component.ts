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
import {Project} from '../../_models/project';
import {DataService} from '../../_services/data.service';
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
  public selectedDeployment?: Deployment;
  public isQualityGatesOnly = false;
  public selectedStage?: string;

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
        const changedDeployments = (this.selectedDeployment && service?.deployments.filter(deployment => deployment.name === this.selectedDeployment?.name)) ?? []; // the context of a deployment may change

        if (paramDeployment) {
          this.selectedDeployment = paramDeployment;
        } else if (changedDeployments.length > 0) {
          let deployment;
          if (changedDeployments.length === 1) {
            deployment = changedDeployments[0];
          } else {
            deployment = changedDeployments.find(d => d.stages.some(s => this.selectedDeployment?.stages.some(sd => s.stageName === sd.stageName)));
          }
          if (deployment) {
            this.selectedDeployment = deployment;
          }
        } else {
          const routeUrl = this.router.createUrlTree(['/project', project.projectName, 'service', params.serviceName]);
          this.location.go(routeUrl.toString());
        }
        this.selectedStage = params.stage;
      }
    });

    project$.subscribe(project => {
      if (this.selectedDeployment) { // the selected deployment gets lost if the project is updated, because the deployments are rebuild
        this.selectedDeployment = project.getServices().find(s => s.serviceName === this.selectedDeployment?.service)?.deployments.find(d => d.shkeptncontext === this.selectedDeployment?.shkeptncontext);
      }
      this.dataService.loadOpenRemediations(project);
      this.project = project;
      this._changeDetectorRef.markForCheck();
    });
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
