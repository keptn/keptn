import { Store } from '@ngrx/store';
import { Injectable } from '@angular/core';
import {
  fromKeptInfo,
  fromLatestSequences,
  fromMetadata,
  fromProjects,
  qualityGatesOnly,
} from './root.store.selectors';
import { loadLatestSequences, loadRootState } from './root.store.actions';
import { State } from './root.store.reducer';

@Injectable({
  providedIn: 'root',
})
export class RootStoreFacade {
  keptInfo$ = this.store.select(fromKeptInfo);
  metadata$ = this.store.select(fromMetadata);
  projects$ = this.store.select(fromProjects);
  qualityGatesOnly$ = this.store.select(qualityGatesOnly);
  latestSequences$ = this.store.select(fromLatestSequences);

  constructor(private store: Store<State>) {}

  loadRootState(): void {
    this.store.dispatch(loadRootState());
  }

  refreshSequences(): void {
    this.store.dispatch(loadLatestSequences());
  }
}
