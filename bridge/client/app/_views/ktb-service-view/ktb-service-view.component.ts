import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  OnInit,
  ViewEncapsulation
} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Subject, timer} from 'rxjs';
import {filter, map, switchMap, takeUntil} from 'rxjs/operators';
import {Project} from '../../_models/project';
import {DataService} from '../../_services/data.service';
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
  private _projectTimerInterval = 30 * 1000;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute) { }

  ngOnInit() {
    this.dataService.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(keptnInfo => {
        this.isQualityGatesOnly = !keptnInfo.bridgeInfo.keptnInstallationType?.includes('CONTINUOUS_DELIVERY');
        this._changeDetectorRef.markForCheck();
      });

    this.dataService._remediationsUpdated
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this._changeDetectorRef.markForCheck();
      });

    const params$ = this.route.params
      .pipe(takeUntil(this.unsubscribe$));

    const project$ = params$.pipe(
      switchMap(params => this.dataService.getProject(params.projectName)),
      takeUntil(this.unsubscribe$)
      );

    const timer$ = params$.pipe(
      switchMap((params) => timer(0, this._projectTimerInterval).pipe(map(() => params.projectName))),
      takeUntil(this.unsubscribe$)
    );

    params$.subscribe(params => {
      this.serviceName ??= params.serviceName;
      this._changeDetectorRef.markForCheck();
    });

    project$.subscribe(project => {
      this.dataService.loadOpenRemediations(project);
      this.project = project;
      this._changeDetectorRef.markForCheck();
    });

    timer$.subscribe(projectName => {
      this.dataService.loadProject(projectName);
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
