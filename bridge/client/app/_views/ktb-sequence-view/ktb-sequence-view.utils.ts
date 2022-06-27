import { Sequence } from '../../_models/sequence';

export enum SequencesState {
  LOAD_UNTIL_ROOT,
  UPDATE,
}

export interface ISequenceStateInfo {
  allSequencesLoaded: boolean;
  sequences: Sequence[];
  state: SequencesState;
}

export interface ISequenceState {
  [projectName: string]: ISequenceStateInfo | undefined;
}

export enum FilterName {
  SERVICE = 'Service',
  STAGE = 'Stage',
  SEQUENCE = 'Sequence',
  STATUS = 'Status',
}

export type FilterType = [
  {
    name: FilterName;
    autocomplete: { name: string; value: string }[];
    showInSidebar: boolean;
  },
  ...{ name: string; value: string }[]
];

export interface ISequenceViewState {
  projectName: string;
  eventId?: string;
  sequenceInfo?: ISequenceStateInfo;
}
