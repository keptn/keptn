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
import { Deployment, DeploymentSelection } from '../../_models/deployment';
import {Location} from '@angular/common';

@Component({
  selector: 'ktb-deployment-list[service]',
  templateUrl: './ktb-deployment-list.component.html',
  styleUrls: ['./ktb-deployment-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbDeploymentListComponent implements OnInit, OnDestroy {
  private _service?: Service;
  private projectName?: string;
  private readonly unsubscribe$ = new Subject<void>();
  public _selectedDeploymentInfo?: DeploymentSelection;
  public dataSource = new DtTableDataSource<Deployment>();
  public loading = false;
  public DeploymentClass = Deployment;

  @Output() selectedDeploymentInfoChange: EventEmitter<DeploymentSelection> = new EventEmitter();

  @Input()
  get service(): Service | undefined {
    return this._service;
  }

  set service(service: Service | undefined) {
    if (this._service !== service) {
      this._service = service;
      this._changeDetectorRef.markForCheck();
    }
  }
  @Input()
  get selectedDeploymentInfo(): DeploymentSelection | undefined {
    return this._selectedDeploymentInfo;
  }
  set selectedDeploymentInfo(deployment: DeploymentSelection | undefined) {
    if (this._selectedDeploymentInfo !== deployment) {
      this._selectedDeploymentInfo = deployment;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(public _changeDetectorRef: ChangeDetectorRef, private route: ActivatedRoute, private dataService: DataService,
              private router: Router, private location: Location) { }

  ngOnInit(): void {
    this.dataService.changedDeployments
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((deployments) => {
        if (deployments.some(d => d.service === this.service?.serviceName)) {
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
      this.projectName = project?.projectName;
      this.updateDataSource();
    });
  }

  private updateDataSource(count = -1): void {
    this.dataSource.data = (count !== -1 ? this.service?.deployments.slice(0, count) : this.service?.deployments) ?? [];
    this._changeDetectorRef.markForCheck();
  }

  public selectDeployment(deployment: Deployment, stageName?: string): void {
    if (this.selectedDeploymentInfo?.deployment.shkeptncontext !== deployment.shkeptncontext || stageName) {
      stageName ??= deployment.stages[deployment.stages.length - 1].stageName;
      const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', deployment.service, 'context', deployment.shkeptncontext, 'stage', stageName]);
      this.location.go(routeUrl.toString());
      this.selectedDeploymentInfo = {deployment, stage: stageName};
      this.selectedDeploymentInfoChange.emit(this.selectedDeploymentInfo);
      this._changeDetectorRef.markForCheck();
    }
  }

  public selectStage(deployment: Deployment, stageName: string, $event: MouseEvent) {
    $event.stopPropagation();
    this.selectDeployment(deployment, stageName);
    this._changeDetectorRef.markForCheck();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
