import { Component, HostBinding, Inject, OnDestroy } from '@angular/core';
import { distinctUntilChanged, filter, map, switchMap, takeUntil } from 'rxjs/operators';
import { combineLatest, Observable, Subject } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { ServiceFilterType } from './ktb-stage-details/ktb-stage-details.component';
import { Stage } from '../../_models/stage';
import { Location } from '@angular/common';

export interface ISelectedStageInfo {
  stage: Stage;
  filterType: ServiceFilterType;
}

@Component({
  selector: 'ktb-environment-view',
  templateUrl: './ktb-environment-view.component.html',
  styleUrls: ['./ktb-environment-view.component.scss'],
  preserveWhitespaces: false,
})
export class KtbEnvironmentViewComponent implements OnDestroy {
  @HostBinding('class') cls = 'ktb-environment-view';
  public project$: Observable<Project>;
  private readonly unsubscribe$ = new Subject<void>();
  public selectedStageInfo?: ISelectedStageInfo;
  public isQualityGatesOnly$: Observable<boolean>;
  public isTriggerSequenceOpen$: Observable<boolean>;

  constructor(
    private dataService: DataService,
    private route: ActivatedRoute,
    private router: Router,
    private location: Location,
    @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number
  ) {
    this.isQualityGatesOnly$ = this.dataService.isQualityGatesOnly;
    this.isTriggerSequenceOpen$ = this.dataService.isTriggerSequenceOpen;
    this.dataService.setIsTriggerSequenceOpen(false);
    const selectedStageName$ = this.route.paramMap.pipe(
      map((params) => params.get('stageName')),
      takeUntil(this.unsubscribe$)
    );

    const paramFilterType$ = this.route.queryParamMap.pipe(map((params) => params.get('filterType')));
    const projectName$ = this.route.paramMap.pipe(
      map((params) => params.get('projectName')),
      filter((projectName): projectName is string => !!projectName),
      distinctUntilChanged()
    );

    this.project$ = projectName$.pipe(
      switchMap((projectName) => this.dataService.getProject(projectName)),
      filter((project): project is Project => !!project?.projectDetailsLoaded)
    );

    combineLatest([selectedStageName$, paramFilterType$, this.project$])
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(([stageName, filterType, project]) => {
        const stage = project.stages.find((s) => s.stageName === stageName);
        if (stage) {
          this.selectedStageInfo = { stage: stage, filterType: (filterType as ServiceFilterType) ?? undefined };
        }
      });

    const projectTimer$ = AppUtils.createTimer(0, initialDelayMillis).pipe(
      switchMap(() => projectName$),
      takeUntil(this.unsubscribe$)
    );

    projectTimer$.subscribe((projectName) => {
      this.dataService.loadProject(projectName);
    });
  }

  public setSelectedStageInfo(projectName: string, stageInfo: ISelectedStageInfo): void {
    this.selectedStageInfo = stageInfo;
    this.setLocation(projectName, stageInfo);
  }

  private setLocation(projectName: string, stageInfo: ISelectedStageInfo): void {
    const url = this.router.createUrlTree(
      ['/project', projectName, 'environment', 'stage', stageInfo.stage.stageName],
      {
        queryParams: { filterType: stageInfo.filterType },
      }
    );
    this.location.go(url.toString());
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
