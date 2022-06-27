import { ChangeDetectionStrategy, Component } from '@angular/core';
import { catchError, filter, map, startWith, switchMap } from 'rxjs/operators';
import { Observable, of, throwError } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { DataService } from '../_services/data.service';
import { environment } from '../../environments/environment';
import { KeptnService } from '../../../shared/models/keptn-service';
import {
  EvaluationBoardParams,
  EvaluationBoardState,
  EvaluationBoardStateLoading,
  EvaluationBoardStatus,
} from './evaluation-board-state';
import { DateUtil } from '../_utils/date.utils';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'ktb-evaluation-board',
  templateUrl: './evaluation-board.component.html',
  styleUrls: ['./evaluation-board.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class EvaluationBoardComponent {
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public hasHistory: boolean;
  public EvaluationBoardState = EvaluationBoardStatus;
  public state$: Observable<EvaluationBoardState>;

  constructor(private location: Location, private route: ActivatedRoute, private dataService: DataService) {
    this.hasHistory = window.history.length > 1;
    this.dataService.setProjectName(''); // else in the app-header the latest projectName will be shown until the traces are loaded

    this.state$ = this.route.paramMap.pipe(
      map((params) => ({ keptnContext: params.get('shkeptncontext'), eventSelector: params.get('eventselector') })),
      filter((params): params is EvaluationBoardParams => !!params.keptnContext),
      switchMap((params) =>
        // Get evaluations for the given keptnContext
        this.dataService
          .getTracesByContext(params.keptnContext, EventTypes.EVALUATION_FINISHED, KeptnService.LIGHTHOUSE_SERVICE)
          .pipe(
            map((evaluations) =>
              params.eventSelector
                ? evaluations.filter((t) => t.id === params.eventSelector || t.data.stage === params.eventSelector)
                : evaluations
            ),
            map((evaluations) => evaluations.sort(DateUtil.compareTraceTimesDesc)),
            map((evaluations) => ({ evaluations, keptnContext: params.keptnContext })),
            catchError(() => throwError(() => ({ keptnContext: params.keptnContext })))
          )
      ),
      switchMap((data) =>
        // Only the triggered events have the previous data => configurationChange.values.image does not exist on finished events
        this.dataService
          .getTracesByContext(data.keptnContext, EventTypes.EVALUATION_TRIGGERED, undefined, undefined, 1)
          .pipe(
            map((root) => {
              return {
                ...data,
                artifact: root[0].getConfigurationChangeImage(),
                deploymentName: root[0].getShortImageName() || root[0].service || '',
              };
            })
          )
      ),
      switchMap((data): Observable<EvaluationBoardState> => {
        const { project: projectName, stage: stageName, service: serviceName } = data.evaluations[0];
        if (!projectName || !stageName || !serviceName) {
          return of({ state: EvaluationBoardStatus.ERROR, kind: 'trace', keptnContext: data.keptnContext });
        }

        this.dataService.setProjectName(projectName);
        // Get the latest deployment context that matches with the deployment contexts in the service screen
        // An evaluation can be triggered without a deployment,
        //   so there is the case that the evaluation keptnContext does not match with any deployment keptnContext
        return this.dataService.getService(projectName, stageName, serviceName).pipe(
          map((service) => service.deploymentContext),
          map((serviceKeptnContext) => ({
            evaluations: data.evaluations,
            deploymentName: data.deploymentName,
            artifact: data.artifact,
            serviceKeptnContext,
            state: EvaluationBoardStatus.LOADED,
          }))
        );
      }),
      catchError((error: HttpErrorResponse | { keptnContext: string }): Observable<EvaluationBoardState> => {
        if ('keptnContext' in error) {
          return of({ state: EvaluationBoardStatus.ERROR, kind: 'trace', keptnContext: error.keptnContext });
        }
        return of({ state: EvaluationBoardStatus.ERROR, kind: 'default' });
      }),
      startWith({ state: EvaluationBoardStatus.LOADING } as EvaluationBoardStateLoading)
    );
  }

  public goBack(): void {
    this.location.back();
  }
}
