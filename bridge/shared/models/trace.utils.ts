import { ICombinedTrace, ITrace } from './trace';
import { EventTypes } from '../interfaces/event-types';
import { ProblemStates } from './problem-states';
import { ResultTypes } from './result-types';
import { ApprovalStates } from './approval-states';

export function getFinishedEvent(trace: ICombinedTrace): ICombinedTrace | undefined {
  return isFinishedEvent(trace) ? trace : trace.traces.find((t) => isFinishedEvent(t));
}

export function isFinishedEvent(trace: ITrace | ICombinedTrace): boolean {
  return trace.type.endsWith('.finished');
}

export function isProblem(trace: ITrace): boolean {
  return trace.type === EventTypes.PROBLEM_DETECTED || trace.type === EventTypes.PROBLEM_OPEN;
}

export function isProblemResolvedOrClosed(trace: ICombinedTrace): boolean {
  if (!trace.traces || trace.traces.length === 0) {
    return trace.data.State === ProblemStates.RESOLVED || trace.data.State === ProblemStates.CLOSED;
  } else {
    return trace.traces.some((t) => isProblem(t) && isProblemResolvedOrClosed(trace));
  }
}

export function isStartedEvent(trace: ITrace): boolean {
  return trace.type.endsWith('.started');
}

export function isFaulty(trace: ICombinedTrace, stageName?: string): boolean {
  let result = false;
  if (trace.data) {
    if (isFailed(trace) || trace.traces.some((t) => isFaulty(t))) {
      result = stageName ? trace.data.stage === stageName : true;
    }
  }
  return result;
}

export function isFailed(trace: ICombinedTrace): boolean {
  return (
    getFinishedEvent(trace)?.data.result === ResultTypes.FAILED || (isApprovalFinished(trace) && isDeclined(trace))
  );
}

export function isApprovalFinished(trace: ITrace): boolean {
  return trace.type === EventTypes.APPROVAL_FINISHED;
}

export function isDeclined(trace: ITrace): boolean {
  return trace.data.approval?.result === ApprovalStates.DECLINED;
}

export function isSuccessfulRemediation(trace: ICombinedTrace): boolean {
  if (!trace.traces || trace.traces.length === 0) {
    return trace.type.endsWith(EventTypes.REMEDIATION_FINISHED_SUFFIX) && trace.data.result !== ResultTypes.FAILED;
  } else {
    return trace.traces.some((t) => isSuccessfulRemediation(t));
  }
}
