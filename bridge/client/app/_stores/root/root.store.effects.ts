import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import {
  keptnInfoLoaded,
  latestSequencesForProjectLoaded,
  loadLatestSequences,
  loadRootState,
  metadataErrored,
  metadataLoaded,
  projectsLoaded,
  refreshProjects,
} from './root.store.actions';
import { catchError, filter, map, mergeMap } from 'rxjs/operators';
import { combineLatest, forkJoin, merge, Observable, of } from 'rxjs';
import { KeptnInfo } from '../../_models/keptn-info';
import { Store } from '@ngrx/store';
import { State } from './root.store.reducer';
import { fromProjects, fromProjectsPageSize } from './root.store.selectors';
import { ApiService } from '../../_services/api.service';
import { IProject } from '../../../../shared/interfaces/project';
import { LoadingState } from '../store';
import { ISequence } from '../../../../shared/interfaces/sequence';

@Injectable()
export class RootStoreEffects {
  loadKeptInfo$ = createEffect(() =>
    this.actions$.pipe(
      ofType(loadRootState),
      mergeMap(() => this.loadKeptnInfo()),
      map((keptnInfo) => keptnInfoLoaded({ keptnInfo }))
    )
  );

  loadMetadata$ = createEffect(() =>
    this.actions$.pipe(
      ofType(loadRootState),
      mergeMap(() => this.apiService.getMetadata().pipe(catchError(() => of(undefined)))),
      map((metadata) => (metadata ? metadataLoaded({ metadata }) : metadataErrored()))
    )
  );

  loadProjects$ = createEffect(() => {
    const pageSize$ = this.store.select(fromProjectsPageSize);
    const loadAction$ = this.actions$.pipe(ofType(keptnInfoLoaded, refreshProjects));
    return combineLatest([pageSize$, loadAction$]).pipe(
      mergeMap(([pageSize]) => this.loadProjects(pageSize)),
      map((projects) => projectsLoaded({ projects }))
    );
  });

  loadLatestSequences$ = createEffect(() => {
    const projects$ = this.store.select(fromProjects).pipe(
      filter((value) => value.call === LoadingState.LOADED),
      map((value) => value.data)
    );
    const loadAction$ = this.actions$.pipe(ofType(loadLatestSequences));
    return combineLatest([loadAction$, projects$]).pipe(
      mergeMap(([, projects]) => merge(...this.loadSequences(projects))),
      map((value) => latestSequencesForProjectLoaded(value))
    );
  });

  constructor(private store: Store<State>, private actions$: Actions, private apiService: ApiService) {}

  loadKeptnInfo(): Observable<KeptnInfo> {
    return this.apiService.getKeptnInfo().pipe(
      mergeMap((bridgeInfo) => {
        const availableVersionsObs = bridgeInfo.enableVersionCheckFeature
          ? this.apiService.getAvailableVersions().pipe(
              catchError(() => {
                return of(undefined);
              })
            )
          : of(undefined);
        return forkJoin([of(bridgeInfo), availableVersionsObs]);
      }),
      map(([bridgeInfo, availableVersions]) => ({ bridgeInfo, availableVersions } as KeptnInfo))
    );
  }

  loadProjects(pageSize?: number): Observable<IProject[]> {
    return this.apiService.getProjects(pageSize || 50).pipe(map((res) => (res ? res.projects : [])));
  }

  private loadSequences(projects: IProject[]): Observable<{ projectName: string; sequences: ISequence[] }>[] {
    return projects.map((project) =>
      this.apiService
        .getSequences(project.projectName, 5)
        .pipe(map((response) => ({ projectName: project.projectName, sequences: response.body?.states ?? [] })))
    );
  }
}
