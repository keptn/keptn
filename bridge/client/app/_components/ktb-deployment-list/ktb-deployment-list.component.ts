import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output
} from '@angular/core';
import {Service} from '../../_models/service';
import {switchMap, takeUntil} from 'rxjs/operators';
import {ActivatedRoute, Router} from '@angular/router';
import {combineLatest, Subject} from 'rxjs';
import {DataService} from '../../_services/data.service';
import {DtTableDataSource} from '@dynatrace/barista-components/table';
import {Deployment} from '../../_models/deployment';
import {Location} from '@angular/common';

@Component({
  selector: 'ktb-deployment-list',
  templateUrl: './ktb-deployment-list.component.html',
  styleUrls: ['./ktb-deployment-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbDeploymentListComponent implements OnInit, OnDestroy {
  private _service: Service;
  private projectName: string;
  private readonly unsubscribe$ = new Subject<void>();
  public _selectedDeployment: Deployment;
  public dataSource = new DtTableDataSource();
  public pageSize: number;
  public minPageSize: number;
  public loading = false;

  @Output() selectedDeploymentChange: EventEmitter<Deployment> = new EventEmitter();

  @Input()
  get service(): Service {
    return this._service;
  }

  set service(service: Service) {
    if (this._service !== service) {
      this._service = service;
      this._changeDetectorRef.markForCheck();
    }
  }
  @Input()
  get selectedDeployment(): Deployment {
    return this._selectedDeployment;
  }
  set selectedDeployment(deployment: Deployment) {
    if (this._selectedDeployment !== deployment) {
      this._selectedDeployment = deployment;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(public _changeDetectorRef: ChangeDetectorRef, private route: ActivatedRoute, private dataService: DataService, private router: Router, private location: Location) { }

  ngOnInit(): void {
    this.dataService.changedDeployments
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((deployments) => {
        if (deployments.some(d => d.service === this.service.serviceName)) {
          this._changeDetectorRef.markForCheck();
        }
      });

    const params$ = this.route.params
      .pipe(takeUntil(this.unsubscribe$));

    const project$ = params$.pipe(
      switchMap(params => this.dataService.getProject(params.projectName)),
      takeUntil(this.unsubscribe$)
    );

    combineLatest([params$, project$]).subscribe(([params, project]) => {
      this.projectName = project.projectName;
      if (params.shkeptncontext && this.service.serviceName === params.serviceName) {
        const paramDeployment = this.service.deployments.find(deployment => deployment.shkeptncontext === params.shkeptncontext);
        const changedDeployment = this.selectedDeployment && this.service.deployments.filter(deployment => deployment.name === this.selectedDeployment.name); // the context of a deployment may change

        if (paramDeployment) {
          this.selectDeployment(paramDeployment, !params.stage);
        } else if (changedDeployment?.length > 0) {
          let deployment;
          if (changedDeployment.length === 1) {
            deployment = changedDeployment[0];
          } else {
            deployment = changedDeployment.find(d => d.stages.some(s => this.selectedDeployment.stages.some(sd => s.stageName === sd.stageName)));
          }
          if (deployment) {
            this.selectDeployment(deployment, true);
          }
        } else {
          const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', params.serviceName]);
          this.location.go(routeUrl.toString());
        }
      }
      this.minPageSize = this.service.deployments.length;
      this.updateDataSource();
    });
  }

  private updateDataSource(count = -1): void {
    this.dataSource.data = count !== -1 ? this.service.deployments.slice(0, count) : this.service.deployments;
    this.pageSize = this.dataSource.data.length;
    this._changeDetectorRef.markForCheck();
  }

  public selectDeployment(deployment: Deployment, redirect = true): void {
    if (this.selectedDeployment?.shkeptncontext !== deployment.shkeptncontext) {
      this.selectedDeployment = deployment;

      if (redirect) {
        const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', deployment.service, 'context', deployment.shkeptncontext]);
        this.location.go(routeUrl.toString());
      }
      this.selectedDeploymentChange.emit(this.selectedDeployment);
      this._changeDetectorRef.markForCheck();
    }
  }

  public selectStage(deployment: Deployment, stageName: string, $event) {
    $event.stopPropagation();
    this.selectDeployment(deployment, false);
    this.router.navigate(['/', 'project', this.projectName, 'service', this.service.serviceName, 'context', deployment.shkeptncontext, 'stage', stageName]);
    this._changeDetectorRef.markForCheck();
  }

  public toDeployment(row: Deployment): Deployment {
    return row;
  }

  loadVersions(): void {
    if (this.pageSize === this.minPageSize) {
      if (this.service.allDeploymentsLoaded) {
        this.updateDataSource();
      } else {
        this.loading = true;
        this._changeDetectorRef.markForCheck();

        this.dataService.getDeploymentsOfService(this.projectName, this.service.serviceName).subscribe(deployments => {
          this.service.deployments = [...this.service.deployments, ...deployments];
          this.service.allDeploymentsLoaded = true;
          this.loading = false;
          this.updateDataSource();
        });
      }
    } else {
      this.updateDataSource(this.minPageSize);
    }
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
