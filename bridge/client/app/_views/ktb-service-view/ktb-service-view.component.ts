import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  OnInit,
  ViewEncapsulation
} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Observable, Subject, timer} from 'rxjs';
import {filter, takeUntil} from 'rxjs/operators';
import {Project} from '../../_models/project';
import {DataService} from '../../_services/data.service';
import {Location} from '@angular/common';
import { Deployment } from 'client/app/_models/deployment';

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
  public project: Project;
  public serviceName: string;
  public selectedDeployment: Deployment;
  public isQualityGatesOnly: boolean;
  private _remediationTimerInterval = 30 * 1000;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute, private router: Router, private location: Location) { }

  ngOnInit() {
    this.dataService.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(keptnInfo => {
        this.isQualityGatesOnly = !keptnInfo.bridgeInfo.keptnInstallationType?.includes('CONTINUOUS_DELIVERY');
      });

    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.serviceName = params.serviceName;

        this.dataService.getProject(params.projectName).subscribe(project => {
          this.project = project;
          this.serviceName = params.serviceName;

          this.dataService._remediationsUpdated
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe(() => {
              this._changeDetectorRef.markForCheck();
          });

          timer(0, this._remediationTimerInterval)
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe(() => {
              this.dataService.loadOpenRemediations(project);
            });
        });
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
  }
}
