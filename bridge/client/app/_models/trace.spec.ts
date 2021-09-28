import { Trace } from './trace';
import { waitForAsync } from '@angular/core/testing';
import { EvaluationTracesMock, RootTracesMock } from './trace.mock';

describe('Trace', () => {
  it(
    'should create instances from json',
    waitForAsync(() => {
      const rootTraces: Trace[] = RootTracesMock.map((trace: unknown) => Trace.fromJSON(trace));
      const evaluationTraces: Trace[] = EvaluationTracesMock.map((trace: unknown) => Trace.fromJSON(trace));

      expect(rootTraces[0] instanceof Trace).toEqual(true);

      expect(rootTraces[0].type).toEqual('sh.keptn.event.service.create.started');
      expect(rootTraces[0].getLabel()).toEqual('create');
      expect(rootTraces[0].getShortImageName()).toEqual(undefined);
      expect(rootTraces[0].getIcon()).toEqual('information');
      expect(rootTraces[0].isFaulty()).toEqual(false);
      expect(rootTraces[0].isWarning()).toEqual(false);
      expect(rootTraces[0].isSuccessful()).toEqual(false);
      expect(rootTraces[0].project).toEqual('sockshop');
      expect(rootTraces[0].service).toEqual('carts');

      expect(rootTraces[1].type).toEqual('sh.keptn.event.dev.artifact-delivery.triggered');
      expect(rootTraces[1].getLabel()).toEqual('artifact-delivery');
      expect(rootTraces[1].getShortImageName()).toEqual('carts:0.10.1');
      expect(rootTraces[1].getIcon()).toEqual('duplicate');
      expect(rootTraces[1].isFaulty()).toEqual(false);
      expect(rootTraces[1].isWarning()).toEqual(false);
      expect(rootTraces[1].isSuccessful()).toEqual(false);
      expect(rootTraces[1].project).toEqual('sockshop');
      expect(rootTraces[1].service).toEqual('carts');

      expect(rootTraces[2].type).toEqual('sh.keptn.event.dev.artifact-delivery.triggered');
      expect(rootTraces[2].getLabel()).toEqual('artifact-delivery');
      expect(rootTraces[2].getShortImageName()).toEqual('carts:0.10.2');
      expect(rootTraces[2].getIcon()).toEqual('duplicate');
      expect(rootTraces[2].isFaulty()).toEqual(false);
      expect(rootTraces[2].isWarning()).toEqual(false);
      expect(rootTraces[2].isSuccessful()).toEqual(false);
      expect(rootTraces[2].project).toEqual('sockshop');
      expect(rootTraces[2].service).toEqual('carts');

      expect(rootTraces[8].type).toEqual('sh.keptn.event.dev.artifact-delivery.triggered');
      expect(rootTraces[8].getLabel()).toEqual('artifact-delivery');
      expect(rootTraces[8].getShortImageName()).toEqual(undefined);
      expect(rootTraces[8].getIcon()).toEqual('duplicate');
      expect(rootTraces[8].isFaulty()).toEqual(false);
      expect(rootTraces[8].isWarning()).toEqual(false);
      expect(rootTraces[8].isSuccessful()).toEqual(false);
      expect(rootTraces[8].project).toEqual('keptn');
      expect(rootTraces[8].service).toEqual('control-plane');

      expect(evaluationTraces[0].type).toEqual('sh.keptn.event.evaluation.finished');
      expect(evaluationTraces[0].getLabel()).toEqual('evaluation');
      expect(evaluationTraces[0].getIcon()).toEqual('traffic-light');
      expect(evaluationTraces[0].isFaulty()).toEqual(false);
      expect(evaluationTraces[0].isWarning()).toEqual(false);
      expect(evaluationTraces[0].isSuccessful(evaluationTraces[0].data.stage)).toEqual(true);
      expect(evaluationTraces[0].project).toEqual('sockshop');
      expect(evaluationTraces[0].service).toEqual('carts');
    })
  );
});
