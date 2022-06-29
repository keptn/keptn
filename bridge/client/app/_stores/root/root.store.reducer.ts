import { createReducer, on } from '@ngrx/store';
import {
  keptnInfoLoaded,
  latestSequencesForProjectLoaded,
  loadRootState,
  metadataErrored,
  metadataLoaded,
  projectsLoaded,
} from './root.store.actions';
import { ApiCall, LoadingState } from '../store';
import { KeptnInfo } from '../../_models/keptn-info';
import { IMetadata } from '../../_interfaces/metadata';
import { IProject } from '../../../../shared/interfaces/project';
import { ISequence } from '../../../../shared/interfaces/sequence';

export interface RootState {
  keptInfo: ApiCall<KeptnInfo | undefined>;
  metadata: ApiCall<IMetadata | undefined>;
  projects: ApiCall<IProject[]>;
  latestSequences: Record<string, ISequence[]>;
}

export const initialRootState: RootState = {
  keptInfo: { call: LoadingState.INIT, data: undefined },
  metadata: { call: LoadingState.INIT, data: undefined },
  projects: { call: LoadingState.INIT, data: [] },
  latestSequences: {},
};

export interface State {
  root: RootState;
}

export const rootStoreReducer = createReducer(
  initialRootState,

  on(loadRootState, (state) => ({
    ...state,
    keptInfo: { ...state.keptInfo, call: LoadingState.LOADING },
    metadata: { ...state.metadata, call: LoadingState.LOADING },
    projects: { ...state.projects, call: LoadingState.LOADING },
  })),

  on(keptnInfoLoaded, (state, { keptnInfo }) => ({
    ...state,
    keptInfo: { data: keptnInfo, call: LoadingState.LOADED },
  })),

  on(metadataLoaded, (state, { metadata }) => ({
    ...state,
    metadata: { data: metadata, call: LoadingState.LOADED },
  })),

  on(metadataErrored, (state) => ({
    ...state,
    metadata: { ...state.metadata, call: { errorMsg: 'error' } },
  })),

  on(projectsLoaded, (state, { projects }) => ({
    ...state,
    projects: { ...state.projects, data: projects, call: LoadingState.LOADED },
  })),

  on(latestSequencesForProjectLoaded, (state, { projectName, sequences }) => ({
    ...state,
    latestSequences: { ...state.latestSequences, ...{ [projectName]: sequences } },
  }))
);
