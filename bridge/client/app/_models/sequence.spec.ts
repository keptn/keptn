import { createSequenceStateInfo, Sequence } from './sequence';
import {
  SequenceResponseMock,
  SequenceResponseWithDevAndStagingMock,
  SequenceResponseWithoutFailing,
} from '../_services/_mockData/sequences.mock';
import { SequenceState } from '../../../shared/interfaces/sequence';
import { EvaluationTraceResponse } from '../_services/_mockData/evaluations.mock';
import { Trace } from './trace';
import { EventState } from '../../../shared/models/event-state';
import { ResultTypes } from '../../../shared/models/result-types';
import { RemediationAction } from '../../../shared/models/remediation-action';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { AppUtils } from '../_utils/app.utils';

const remediationAction = {
  name: 'Failure rate increase',
  state: EventState.FINISHED,
  action: 'my action',
  result: ResultTypes.PASSED,
  description: '',
};
const devTraceObj = {
  data: {
    project: 'sockshop',
    service: 'carts',
    stage: 'dev',
    labels: {
      label1: 'label1',
    },
  },
  id: 'id1',
  shkeptncontext: 'keptnContext',
  time: '',
  type: EventTypes.DEPLOYMENT_FINISHED,
};

const devTrace = Trace.fromJSON(devTraceObj);

const stagingTraceObj = {
  data: {
    project: 'sockshop',
    service: 'carts',
    stage: 'staging',
    labels: {
      label2: 'label2',
    },
  },
  id: 'id2',
  shkeptncontext: 'keptnContext',
  time: '',
  type: EventTypes.DEPLOYMENT_FINISHED,
};
const stagingTrace = Trace.fromJSON(stagingTraceObj);

describe('Sequence', () => {
  it('should correctly create new class', () => {
    const sequence = getDefaultSequence();
    expect(sequence).toBeInstanceOf(Sequence);
  });

  it('should correctly create new class with extended properties', () => {
    const sequence = Sequence.fromJSON(getSequenceObjectWithEvaluationAndRemediation());
    expect(sequence.stages[0].latestEvaluationTrace).toBeInstanceOf(Trace);
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    expect(sequence.stages[0].actions[0]).toBeInstanceOf(RemediationAction);
  });

  it('should return short type delivery', () => {
    expect(Sequence.getShortType('sh.keptn.event.dev.delivery.finished')).toBe('delivery');
    expect(Sequence.getShortType('sh.keptn.event.delivery.finished')).toBe('delivery');
  });

  it('should return full type if type is invalid', () => {
    expect(Sequence.getShortType('sh.keptn.event.dev.delivery.finished.finished')).toBe(
      'sh.keptn.event.dev.delivery.finished.finished'
    );
  });

  it('should be finished if it is aborted, finished, succeeded or timed out', () => {
    const sequence = getDefaultSequence();
    for (const state of [
      SequenceState.ABORTED,
      SequenceState.FINISHED,
      SequenceState.TIMEDOUT,
      SequenceState.SUCCEEDED,
    ]) {
      sequence.state = state;
      const actual = createSequenceStateInfo(sequence);
      expect(sequence.isFinished()).toBe(true);
      expect(actual.finished).toBe(true);
    }
  });

  it('should not be finished', () => {
    const sequence = getDefaultSequence();
    const { ABORTED, FINISHED, TIMEDOUT, SUCCEEDED, ...states } = SequenceState;
    for (const state of Object.values(states)) {
      sequence.state = state;
      const actual = createSequenceStateInfo(sequence);
      expect(sequence.isFinished()).toBe(false);
      expect(actual.finished).toBe(false);
    }
  });

  it('should be finished stage if it is aborted, finished, succeeded or timed out', () => {
    const sequence = getDefaultSequence();
    for (const state of [
      SequenceState.ABORTED,
      SequenceState.FINISHED,
      SequenceState.TIMEDOUT,
      SequenceState.SUCCEEDED,
    ]) {
      sequence.stages[0].state = state;
      const actual = createSequenceStateInfo(sequence, 'dev');
      expect(sequence.isFinished('dev')).toBe(true);
      expect(actual.finished).toBe(true);
    }
  });

  it('should not be finished stage', () => {
    const sequence = getDefaultSequence();
    const { ABORTED, FINISHED, TIMEDOUT, SUCCEEDED, ...states } = SequenceState;
    for (const state of Object.values(states)) {
      sequence.stages[0].state = state;
      const actual = createSequenceStateInfo(sequence, 'dev');
      expect(sequence.isFinished('dev')).toBe(false);
      expect(actual.finished).toBe(false);
    }
  });

  it('should return all stage names', () => {
    const sequence = getSequenceWithTwoStages();
    expect(sequence.getStages()).toEqual(['dev', 'staging']);
  });

  it('should return last stage', () => {
    const sequence = getSequenceWithTwoStages();
    expect(sequence.getLastStage()).toBe('staging');
  });

  it('should be faulty if it has failed event', () => {
    const sequence = getSequenceWithTwoStages();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isFaulty()).toBe(true);
    expect(actual.faulty).toBe(true);
  });

  it('should have faulty stage if it has failed event', () => {
    const sequence = getSequenceWithTwoStages();
    const actual = createSequenceStateInfo(sequence, 'staging');
    expect(sequence.isFaulty('staging')).toBe(true);
    expect(actual.faulty).toBe(true);
  });

  it('should not be faulty if it does not have failed event', () => {
    const sequence = getEvaluationSequenceWithoutFailing();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isFaulty()).toBe(false);
    expect(actual.faulty).toBe(false);
  });

  it('should not have faulty stage if it has failed event', () => {
    const sequence = getSequenceWithTwoStages();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isFaulty('dev')).toBe(false);
    expect(actual.faulty).toBe(false);
  });

  it('should be faulty if it timed out', () => {
    const sequence = getSequenceWithTwoStages();
    sequence.state = SequenceState.TIMEDOUT;
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isFaulty()).toBe(true);
    expect(actual.faulty).toBe(true);
  });

  it('should have faulty stage if it timed out', () => {
    const sequence = getSequenceWithTwoStages();
    sequence.stages[0].state = SequenceState.TIMEDOUT;
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isFaulty('dev')).toBe(true);
    expect(actual.faulty).toBe(true);
  });

  it('should return evaluation of specific stage', () => {
    const sequence = getEvaluationSequenceWithoutFailing();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.getEvaluation('dev')).toEqual({
      result: ResultTypes.PASSED,
      score: 0,
    });
    expect(actual.evaluation).toEqual({
      result: ResultTypes.PASSED,
      score: 0,
    });
  });

  it('should not return evaluation of not existing stage', () => {
    const sequence = getEvaluationSequenceWithoutFailing();
    const actual = createSequenceStateInfo(sequence, 'staging');
    expect(sequence.getEvaluation('staging')).toBeUndefined();
    expect(actual.evaluation).toBeUndefined();
  });

  it('should not return evaluation if latest evaluation does not exist', () => {
    const sequence = getEvaluationSequenceWithoutFailing();
    sequence.stages[0].latestEvaluation = undefined;
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.getEvaluation('dev')).toBeUndefined();
    expect(actual.evaluation).toBeUndefined();
  });

  it('should return evaluation trace of specific stage', () => {
    const sequence = getSequenceWithEvaluationAndRemediation();
    expect(sequence.getEvaluationTrace('dev')).toEqual(EvaluationTraceResponse.data.evaluationHistory[0]);
  });

  it('should not return evaluation trace of not existing stage', () => {
    const sequence = getSequenceWithEvaluationAndRemediation();
    expect(sequence.getEvaluationTrace('staging')).toBeUndefined();
  });

  it('should not return evaluation trace if latest evaluation does not exist', () => {
    const sequence = getSequenceWithEvaluationAndRemediation();
    sequence.stages[0].latestEvaluationTrace = undefined;
    expect(sequence.getEvaluationTrace('dev')).toBeUndefined();
  });

  it('should have pending approval', () => {
    const sequence = getSequenceWithApproval();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.hasPendingApproval()).toBe(true);
    expect(actual.pendingApproval).toBe(true);
  });

  it('should have pending approval for stage', () => {
    const sequence = getSequenceWithApproval();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.hasPendingApproval('dev')).toBe(true);
    expect(actual.pendingApproval).toBe(true);
  });

  it('should not have pending approval', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.hasPendingApproval()).toBe(false);
    expect(actual.pendingApproval).toBe(false);
  });

  it('should not have pending approval for stage', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.hasPendingApproval('dev')).toBe(false);
    expect(actual.pendingApproval).toBe(false);
  });

  it('should have status failed', () => {
    const sequence = getFailedSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getStatus()).toBe('failed');
    expect(actual.statusText).toBe('failed');
  });

  it('should have status succeeded', () => {
    const sequence = getSucceededSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getStatus()).toBe('succeeded');
    expect(actual.statusText).toBe('succeeded');
  });

  it('should have status aborted', () => {
    const sequence = getAbortedSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getStatus()).toBe('aborted');
    expect(actual.statusText).toBe('aborted');
  });

  it('should have status timed out', () => {
    const sequence = getTimedOutSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getStatus()).toBe('timed out');
    expect(actual.statusText).toBe('timed out');
  });

  it('should have status waiting', () => {
    const sequence = getWaitingSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getStatus()).toBe('waiting');
    expect(actual.statusText).toBe('waiting');
  });

  it('should have status fallback to state', () => {
    const sequence = getUnknownSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getStatus()).toBe('unknown');
    expect(actual.statusText).toBe('unknown');
  });

  it('should be loading', () => {
    const sequence = getLoadingSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isLoading()).toBe(true);
    expect(actual.loading).toBe(true);
  });

  it('should be loading in stage', () => {
    const sequence = getLoadingSequence();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isLoading('dev')).toBe(true);
    expect(actual.loading).toBe(true);
  });

  it('should not be loading', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isLoading()).toBe(false);
    expect(actual.loading).toBe(false);
  });

  it('should not be loading in stage', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isLoading('dev')).toBe(false);
    expect(actual.loading).toBe(false);
  });

  it('should be successful', () => {
    const sequence = getSucceededSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isSuccessful()).toBe(true);
    expect(actual.successful).toBe(true);
  });

  it('should be successful in stage', () => {
    const sequence = getSucceededSequence();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isSuccessful('dev')).toBe(true);
    expect(actual.successful).toBe(true);
  });

  it('should not be successful', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isSuccessful()).toBe(false);
    expect(actual.successful).toBe(false);
  });

  it('should not be successful if stage succeeded but has failed', () => {
    const sequence = getSucceededFailedSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isSuccessful()).toBe(false);
    expect(actual.successful).toBe(false);
  });

  it('should not be successful in stage', () => {
    const sequence = getLoadingSequence();
    const actual = createSequenceStateInfo(sequence, 'staging');
    expect(sequence.isSuccessful('staging')).toBe(false);
    expect(actual.successful).toBe(false);
  });

  it('should be warning', () => {
    const sequence = getEvaluationSequenceWarning();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isWarning()).toBe(true);
    expect(actual.warning).toBe(true);
  });

  it('should be warning in stage', () => {
    const sequence = getEvaluationSequenceWarning();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isWarning('dev')).toBe(true);
    expect(actual.warning).toBe(true);
  });

  it('should not be warning if it is failed', () => {
    const sequence = getEvaluationSequenceWarningFailed();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isWarning()).toBe(false);
    expect(actual.warning).toBe(false);
  });

  it('should not be warning in stage if it is failed', () => {
    const sequence = getEvaluationSequenceWarningFailed();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isWarning('dev')).toBe(false);
    expect(actual.warning).toBe(false);
  });

  it('should be remediation', () => {
    const sequence = getDefaultSequence();
    sequence.name = 'remediation';
    expect(sequence.isRemediation()).toBe(true);
  });

  it('should not be remediation', () => {
    const sequence = getDefaultSequence();
    expect(sequence.isRemediation()).toBe(false);
  });

  it('should be paused', () => {
    const sequence = getPausedSequence();
    expect(sequence.isPaused()).toBe(true);
  });

  it('should not be paused', () => {
    const sequence = getDefaultSequence();
    expect(sequence.isPaused()).toBe(false);
  });

  it('should be unknown state', () => {
    const sequence = getDefaultSequence();
    sequence.state = SequenceState.UNKNOWN;
    expect(sequence.isUnknownState()).toBe(true);
  });

  it('should not be unknown state', () => {
    const sequence = getDefaultSequence();
    expect(sequence.isUnknownState()).toBe(false);
  });

  it('should be aborted', () => {
    const sequence = getAbortedSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isAborted()).toBe(true);
    expect(actual.aborted).toBe(true);
  });

  it('should be aborted stage', () => {
    const sequence = getAbortedSequence();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isAborted('dev')).toBe(true);
    expect(actual.aborted).toBe(true);
  });

  it('should not be aborted', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isAborted()).toBe(false);
    expect(actual.aborted).toBe(false);
  });

  it('should not be aborted stage', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isAborted('dev')).toBe(false);
    expect(actual.aborted).toBe(false);
  });

  it('should be timed out', () => {
    const sequence = getTimedOutSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isTimedOut()).toBe(true);
    expect(actual.timedOut).toBe(true);
  });

  it('should be timed out stage', () => {
    const sequence = getTimedOutSequence();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isTimedOut('dev')).toBe(true);
    expect(actual.timedOut).toBe(true);
  });

  it('should not be timed out', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.isTimedOut()).toBe(false);
    expect(actual.timedOut).toBe(false);
  });

  it('should not be timed out stage', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.isTimedOut('dev')).toBe(false);
    expect(actual.timedOut).toBe(false);
  });

  it('should return latest event', () => {
    const sequence = getSequenceWithTwoStages();
    expect(sequence.getLatestEvent()).toEqual({
      type: 'sh.keptn.event.staging.rollback.finished',
      id: 'b05b8f69-4854-46cd-82d7-69ce3ee73652',
      time: '2021-07-15T15:27:14.208Z',
    });
  });

  it('should not return latest event', () => {
    const sequence = getTriggeredSequence();
    expect(sequence.getLatestEvent()).toBeUndefined();
  });

  it('should return icon of delivery sequence', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getIcon()).toBe('duplicate');
    expect(actual.icon).toBe('duplicate');
  });

  it('should return default icon if not found', () => {
    const sequence = getSequenceWithTwoStages();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getIcon()).toBe('information');
    expect(actual.icon).toBe('information');
  });

  it('should return icon of sequence stage', () => {
    const sequence = getSequenceWithTwoStages();
    const actual = createSequenceStateInfo(sequence, 'dev');
    expect(sequence.getIcon('dev')).toBe('duplicate');
    expect(actual.icon).toBe('duplicate');
  });

  it('should return default icon of sequence stage', () => {
    const sequence = getSequenceWithTwoStages();
    const actual = createSequenceStateInfo(sequence, 'staging');
    expect(sequence.getIcon('staging')).toBe('information');
    expect(actual.icon).toBe('information');
  });

  it('should return default icon of sequence stage if latest stage is undefined', () => {
    const sequence = getTriggeredSequence();
    const actual = createSequenceStateInfo(sequence, 'staging');
    expect(sequence.getIcon('staging')).toBe('information');
    expect(actual.icon).toBe('information');
  });

  it('should return default icon of sequence if latest event is undefined', () => {
    const sequence = getTriggeredSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getIcon()).toBe('information');
    expect(actual.icon).toBe('information');
  });

  it('should return pause icon if sequence is paused', () => {
    const sequence = getPausedSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(sequence.getIcon()).toBe('pause');
    expect(actual.icon).toBe('pause');
  });

  it('should return short image', () => {
    const sequence = getDefaultSequence();
    expect(sequence.getShortImageName()).toBe('carts:0.12.1');
  });

  it('should not return short image', () => {
    const sequence = getDefaultSequence();
    sequence.stages[0].image = undefined;
    expect(sequence.getShortImageName()).toBeUndefined();
  });

  it('should return traces of stage', () => {
    const sequence = getSequenceWithTraces();
    expect(sequence.getTraces('dev')).toEqual([devTrace]);
    expect(sequence.getTraces('staging')).toEqual([stagingTrace]);
  });

  it('should not return any traces of stage', () => {
    const sequence = getSequenceWithTraces();
    expect(sequence.getTraces('production')).toEqual([]);
  });

  it('should return first trace that match condition', () => {
    const sequence = getSequenceWithTraces();
    expect(sequence.findTrace((t) => t.data.project === 'sockshop')).toEqual(devTrace);
  });

  it('should not return traces for given condition', () => {
    const sequence = getSequenceWithTraces();
    expect(sequence.findTrace((t) => t.id === 'idNotExisting')).toBeUndefined();
  });

  it('should return last trace that match condition', () => {
    const sequence = getSequenceWithTraces();
    expect(sequence.findLastTrace((t) => t.data.project === 'sockshop')).toEqual(stagingTrace);
  });

  it('should not return last trace for given condition', () => {
    const sequence = getSequenceWithTraces();
    expect(sequence.findLastTrace((t) => t.id === 'idNotExisting')).toBeUndefined();
  });

  it('should return labels of latest trace', () => {
    const sequence = getSequenceWithTraces();
    const map = new Map();
    map.set('label2', 'label2');
    expect(sequence.getLabels()).toEqual(map);
  });

  it('should return labels of first trace', () => {
    const sequence = getSequenceWithTraces();
    const map = new Map();
    sequence.traces[1].data.labels = undefined;
    map.set('label1', 'label1');
    expect(sequence.getLabels()).toEqual(map);
  });

  it('should not return labels', () => {
    const sequence = getSequenceWithTraces();
    sequence.traces[1].data.labels = undefined;
    sequence.traces[0].data.labels = undefined;
    expect(sequence.getLabels()).toBeUndefined();
  });

  it('should return empty remediation actions', () => {
    const sequence = getDefaultSequence();
    expect(sequence.getRemediationActions()).toEqual([]);
  });

  it('should set state', () => {
    const sequence = getDefaultSequence();
    sequence.setState(SequenceState.TIMEDOUT);
    expect(sequence.state).toBe(SequenceState.TIMEDOUT);
  });

  it('should return the start time of a stage', () => {
    const sequence = getSequenceWithMultipleTracesPerStage();
    expect(sequence.getStageTime('dev')).toBe('2022-03-02T12:46:50.991Z');
    expect(sequence.getStageTime('staging')).toBe('2022-03-02T12:55:50.991Z');
  });

  it('should return undefined for time if traces are not loaded', () => {
    const sequence = getDefaultSequence();
    expect(sequence.getStageTime('dev')).toBeUndefined();
  });

  it('should return correct sequence icon when calling createSequenceStateInfo', () => {
    const sequence = getDefaultSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(actual.icon).toBe('duplicate');
  });

  it('should return correct sequence icon when calling createSequenceStateInfo II', () => {
    const sequence = getPausedSequence();
    const actual = createSequenceStateInfo(sequence);
    expect(actual.icon).toBe('pause');
  });

  function getDefaultSequence(): Sequence {
    return Sequence.fromJSON(AppUtils.copyObject(SequenceResponseMock[1]));
  }

  function getSequenceWithTwoStages(): Sequence {
    return Sequence.fromJSON(SequenceResponseWithDevAndStagingMock);
  }

  function getEvaluationSequenceWithoutFailing(): Sequence {
    return Sequence.fromJSON(SequenceResponseWithoutFailing);
  }

  function getSequenceWithEvaluationAndRemediation(): Sequence {
    return Sequence.fromJSON(getSequenceObjectWithEvaluationAndRemediation());
  }

  function getSequenceWithTraces(): Sequence {
    const sequence = Sequence.fromJSON(AppUtils.copyObject(SequenceResponseMock[0]));
    sequence.traces = [
      Trace.fromJSON(AppUtils.copyObject(devTraceObj)),
      Trace.fromJSON(AppUtils.copyObject(stagingTraceObj)),
    ];
    return sequence;
  }

  function getSequenceWithMultipleTracesPerStage(): Sequence {
    const sequence = Sequence.fromJSON(AppUtils.copyObject(SequenceResponseMock[0]));
    sequence.traces = [
      Trace.fromJSON({
        data: {
          project: 'sockshop',
          service: 'carts',
          stage: 'dev',
          labels: {
            label1: 'label1',
          },
        },
        id: 'id1',
        shkeptncontext: 'keptnContext',
        time: '2022-03-02T12:46:50.991Z',
        type: EventTypes.DEPLOYMENT_FINISHED,
      }),
      Trace.fromJSON({
        data: {
          project: 'sockshop',
          service: 'carts',
          stage: 'dev',
          labels: {
            label1: 'label1',
          },
        },
        id: 'id2',
        shkeptncontext: 'keptnContext',
        time: '2022-03-02T12:50:50.991Z',
        type: EventTypes.DEPLOYMENT_FINISHED,
      }),
      Trace.fromJSON({
        data: {
          project: 'sockshop',
          service: 'carts',
          stage: 'staging',
          labels: {
            label1: 'label1',
          },
        },
        id: 'id3',
        shkeptncontext: 'keptnContext',
        time: '2022-03-02T12:55:50.991Z',
        type: EventTypes.DEPLOYMENT_FINISHED,
      }),
    ];
    return sequence;
  }

  function getEvaluationSequenceWarning(): Sequence {
    const sequence = {
      ...SequenceResponseWithoutFailing,
      stages: [
        {
          ...SequenceResponseWithoutFailing.stages[0],
          latestEvaluation: {
            result: ResultTypes.WARNING,
            score: 75,
          },
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getTriggeredSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      stages: [],
    };
    return Sequence.fromJSON(sequence);
  }

  function getEvaluationSequenceWarningFailed(): Sequence {
    const sequence = {
      ...SequenceResponseWithoutFailing,
      stages: [
        {
          ...SequenceResponseWithoutFailing.stages[0],
          latestEvaluation: {
            result: ResultTypes.WARNING,
            score: 75,
          },
          latestFailedEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_FINISHED,
          },
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getSequenceWithApproval(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.APPROVAL_STARTED,
          },
          latestFailedEvent: undefined,
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getPausedSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: SequenceState.PAUSED,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.APPROVAL_STARTED,
          },
          latestFailedEvent: undefined,
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getFailedSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: SequenceState.FINISHED,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          state: SequenceState.FINISHED,
          latestFailedEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_FINISHED,
          },
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getSucceededSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: SequenceState.FINISHED,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          state: SequenceState.SUCCEEDED,
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_FINISHED,
          },
          latestFailedEvent: undefined,
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getSucceededFailedSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: SequenceState.FINISHED,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          state: SequenceState.SUCCEEDED,
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_FINISHED,
          },
          latestFailedEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_FINISHED,
          },
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getWaitingSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: SequenceState.WAITING,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          state: SequenceState.WAITING,
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_TRIGGERED,
          },
          latestFailedEvent: undefined,
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getLoadingSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: SequenceState.STARTED,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          state: SequenceState.STARTED,
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_STARTED,
          },
          latestFailedEvent: undefined,
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getUnknownSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: 'unknown' as SequenceState,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          state: 'unknown' as SequenceState,
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_TRIGGERED,
          },
          latestFailedEvent: undefined,
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getTimedOutSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: SequenceState.TIMEDOUT,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          state: SequenceState.TIMEDOUT,
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_FINISHED,
          },
          latestFailedEvent: undefined,
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getAbortedSequence(): Sequence {
    const sequence = {
      ...SequenceResponseMock[0],
      state: SequenceState.ABORTED,
      stages: [
        {
          ...SequenceResponseMock[0].stages[0],
          state: SequenceState.ABORTED,
          latestEvent: {
            id: 'my Id',
            time: '',
            type: EventTypes.DEPLOYMENT_FINISHED,
          },
        },
      ],
    };
    return Sequence.fromJSON(sequence);
  }

  function getSequenceObjectWithEvaluationAndRemediation(): Sequence {
    const seq = <Sequence>AppUtils.copyObject(SequenceResponseMock[0]);
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    seq.stages[0].latestEvaluationTrace = EvaluationTraceResponse.data.evaluationHistory[0];
    seq.stages[0].actions = [
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      remediationAction,
    ];
    return seq;
  }
});
