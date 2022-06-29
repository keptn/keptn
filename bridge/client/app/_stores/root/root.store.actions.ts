import { createAction, props } from '@ngrx/store';
import { KeptnInfo } from '../../_models/keptn-info';
import { IMetadata } from '../../_interfaces/metadata';
import { IProject } from '../../../../shared/interfaces/project';
import { ISequence } from '../../../../shared/interfaces/sequence';

export const loadRootState = createAction('[Root] Load Root State');
export const keptnInfoLoaded = createAction('[Root] KeptnInfo Loaded', props<{ keptnInfo: KeptnInfo }>());
export const metadataLoaded = createAction('[Root] Metadata Loaded', props<{ metadata: IMetadata }>());
export const metadataErrored = createAction('[Root] Metadata Errored');
export const projectsLoaded = createAction('[Root] Projects Loaded', props<{ projects: IProject[] }>());
export const loadLatestSequences = createAction('[Root] Load Latest Sequences');
export const latestSequencesForProjectLoaded = createAction(
  '[Root] Latest Sequences for Project Loaded',
  props<{ projectName: string; sequences: ISequence[] }>()
);
