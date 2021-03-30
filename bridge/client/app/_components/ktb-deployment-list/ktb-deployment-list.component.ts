import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Service} from '../../_models/service';
import {takeUntil} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {Subject} from 'rxjs';
import {DataService} from '../../_services/data.service';
import {DtTableDataSource} from '@dynatrace/barista-components/table';

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
  public dataSource = new DtTableDataSource();
  public pageSize: number;
  public minPageSize: number;
  public loading = false;
  public gitRemoteURI: string;

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

  constructor(public _changeDetectorRef: ChangeDetectorRef, private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit(): void {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {

        this.dataService.getProject(params.projectName)
          .subscribe(project => {
            this.projectName = project.projectName;
            this.gitRemoteURI = project.gitRemoteURI;
            this.service.deployments = project.getDeploymentsOfService(this.service.serviceName);
            this.minPageSize = this.service.deployments.length;
            this.updateDataSource();
          });

        this.dataService.roots
          .pipe(takeUntil(this.unsubscribe$))
          .subscribe(() => {
            this._changeDetectorRef.markForCheck();
          });
      });
  }

  updateDataSource(count = -1): void {
    this.dataSource.data = count !== -1 ? this.service.deployments.slice(0, count) : this.service.deployments;
    this.pageSize = this.dataSource.data.length;
    this._changeDetectorRef.markForCheck();
  }

  loadVersions(serviceName: string): void {
    if (this.pageSize === this.minPageSize) {
      if (this.service.allDeploymentsLoaded) {
        this.updateDataSource();
      } else {
        this.loading = true;
        this._changeDetectorRef.markForCheck();

        this.dataService.getDeploymentsOfService(this.projectName, serviceName).subscribe(deployments => {
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
