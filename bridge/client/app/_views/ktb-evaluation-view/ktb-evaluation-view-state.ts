import { Trace } from '../../_models/trace';

export enum EvaluationBoardStatus {
  LOADING,
  ERROR,
  LOADED,
}

export type EvaluationBoardParams = { keptnContext: string; eventSelector: string | null };

interface EvaluationBoardStateLoaded {
  evaluations: Trace[];
  serviceKeptnContext?: string;
  state: EvaluationBoardStatus.LOADED;
  artifact?: string;
  deploymentName: string;
}

export interface EvaluationBoardStateLoading {
  state: EvaluationBoardStatus.LOADING;
}

interface EvaluationBoardStateErrorDefault {
  kind: 'default';
}

interface EvaluationBoardStateErrorTrace {
  state: EvaluationBoardStatus.ERROR;
  kind: 'trace';
  keptnContext: string;
}

type EvaluationBoardStateError = (EvaluationBoardStateErrorDefault | EvaluationBoardStateErrorTrace) & {
  state: EvaluationBoardStatus.ERROR;
};

export type KtbEvaluationViewState =
  | EvaluationBoardStateLoaded
  | EvaluationBoardStateLoading
  | EvaluationBoardStateError;
