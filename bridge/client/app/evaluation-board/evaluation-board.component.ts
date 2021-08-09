import {ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {filter, map, takeUntil} from 'rxjs/operators';
import {Subject} from 'rxjs';
import {ActivatedRoute} from '@angular/router';
import {Location} from '@angular/common';
import {Root} from '../_models/root';
import {Trace} from '../_models/trace';
import {ApiService} from '../_services/api.service';
import {EventTypes} from '../../../shared/interfaces/event-types';
import {DataService} from '../_services/data.service';
import {Deployment} from '../_models/deployment';
import {environment} from '../../environments/environment';
import { DateUtil } from '../_utils/date.utils';
import { Project } from '../_models/project';

@Component({
  selector: 'ktb-evaluation-board',
  templateUrl: './evaluation-board.component.html',
  styleUrls: ['./evaluation-board.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class EvaluationBoardComponent implements OnInit, OnDestroy {
  private unsubscribe$ = new Subject<void>();
  private deployments: Deployment[] = [];
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public error?: string;
  public contextId?: string;
  public root?: Root;
  public evaluations?: Trace[];
  public hasHistory: boolean;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private location: Location, private route: ActivatedRoute,
              private apiService: ApiService, private dataService: DataService) {
    this.hasHistory = window.history.length > 1;
    this.dataService.setProjectName(''); // else in the app-header the latest projectName will be shown until the traces are loaded
  }

  ngOnInit() {
    this.route.params
      .pipe(
        takeUntil(this.unsubscribe$),
        // tslint:disable-next-line:no-any
        filter((params: any):
        params is {shkeptncontext: string, eventselector: string | undefined} => !!params.shkeptncontext)
      )
      .subscribe(params => {
        this.contextId = params.shkeptncontext;
        this.apiService.getTraces(this.contextId)
          .pipe(
            map(response => response.body?.events || []),
            map(traces => traces.map(trace => Trace.fromJSON(trace)).sort(DateUtil.compareTraceTimesDesc)),
            takeUntil(this.unsubscribe$)
          ).subscribe((traces: Trace[]) => {
            if (traces.length > 0) {
              this.root = Root.fromJSON(traces[0]);
              this.root.traces = traces;
              this.evaluations = traces.filter(t => t.type === EventTypes.EVALUATION_FINISHED
                                  && (!params.eventselector || t.id === params.eventselector || t.data.stage === params.eventselector));
              if (this.root.project) {
                this.dataService.setProjectName(this.root.project);
                if (this.root.service) {
                  this.setDeployments(this.root.project, this.root.service);
                }
              }
            } else {
              this.error = 'contextError';
              this._changeDetectorRef.markForCheck();
            }
          }, () => {
            this.error = 'error';
            this._changeDetectorRef.markForCheck();
          });
      });
  }

  private setDeployments(projectName: string, serviceName: string): void {
    this.dataService.getProject(projectName)
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((project: Project | undefined): project is Project => !!project)
      )
      .subscribe(project => {
        this.deployments = project.getService(serviceName)?.deployments ?? [];
        this._changeDetectorRef.markForCheck();
      });
  }

  public getDeployment(stage: string) {
    return this.deployments.find(deployment => deployment.stages.find(s => s.stageName === stage));
  }

  public getServiceDetailsLink(shkeptncontext: string, stage: string | undefined): string[] {
    return this.root?.project && this.root?.service && stage
      ? ['/project', this.root.project, 'service', this.root.service, 'context', shkeptncontext, 'stage', stage]
      : [];
  }

  goBack(): void {
    this.location.back();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
