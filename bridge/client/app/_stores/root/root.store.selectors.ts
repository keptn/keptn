import { createFeatureSelector, createSelector } from '@ngrx/store';
import { Features } from '../store';
import { RootState } from './root.store.reducer';

const fromRoot = createFeatureSelector<RootState>(Features.ROOT);

export const fromKeptInfo = createSelector(fromRoot, (state) => state.keptInfo);

export const fromMetadata = createSelector(fromRoot, (state) => state.metadata);

export const fromProjectsPageSize = createSelector(
  fromRoot,
  (state) => state.keptInfo.data?.bridgeInfo.projectsPageSize
);
export const fromProjects = createSelector(fromRoot, (state) => state.projects);

export const qualityGatesOnly = createSelector(fromRoot, (state) => {
  return (
    state.keptInfo.data?.bridgeInfo.showApiToken &&
    !state.keptInfo.data?.bridgeInfo.keptnInstallationType?.includes('CONTINUOUS_DELIVERY')
  );
});

export const fromLatestSequences = createSelector(fromRoot, (state) => state.latestSequences);
