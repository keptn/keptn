import { TestBed } from '@angular/core/testing';
import { EvaluationBoardComponent } from './evaluation-board.component';
import { DataService } from '../_services/data.service';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { BehaviorSubject, firstValueFrom, of, throwError } from 'rxjs';
import { Trace } from '../_models/trace';
import { ActivatedRoute, convertToParamMap, ParamMap } from '@angular/router';
import { Service } from '../_models/service';
import { IService } from '../../../shared/interfaces/service';
import { EvaluationBoardStatus } from './evaluation-board-state';
import { Location } from '@angular/common';
import { skip } from 'rxjs/operators';
import { HttpClientModule } from '@angular/common/http';

describe('ProjectBoardComponent', () => {
  let component: EvaluationBoardComponent;
  let dataService: DataService;
  const paramMap: BehaviorSubject<ParamMap> = new BehaviorSubject<ParamMap>(convertToParamMap({}));
  let location: Location;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [HttpClientModule],
    }).compileComponents();

    dataService = TestBed.inject(DataService);
    location = TestBed.inject(Location);
  });

  it('should return right state', async () => {
    mockComponentData(getDefaultTraces(), getDefaultRootTrace(), getDefaultService(), {
      shkeptncontext: 'myContext',
      eventselector: 'dev',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.LOADED,
      artifact: 'myImage:0.0.1',
      deploymentName: 'myImage:0.0.1',
      serviceKeptnContext: 'myContext',
      evaluations: getDefaultTraces().slice(1).reverse(),
    });
  });

  it('should return and show service name instead of artifact', async () => {
    const rootTrace = getDefaultRootTrace();
    rootTrace.data.configurationChange = undefined;
    mockComponentData(getDefaultTraces(), rootTrace, getDefaultService(), {
      shkeptncontext: 'myContext',
      eventselector: 'dev',
    });

    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.LOADED,
      artifact: undefined,
      deploymentName: 'myService',
      serviceKeptnContext: 'myContext',
      evaluations: getDefaultTraces().slice(1).reverse(),
    });
  });

  it('should should not filter any events', async () => {
    mockComponentData(getDefaultTraces(), getDefaultRootTrace(), getDefaultService(), {
      shkeptncontext: 'myContext',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.LOADED,
      artifact: 'myImage:0.0.1',
      deploymentName: 'myImage:0.0.1',
      serviceKeptnContext: 'myContext',
      evaluations: getDefaultTraces().reverse(),
    });
  });

  it('should filter events by id', async () => {
    mockComponentData(getDefaultTraces(), getDefaultRootTrace(), getDefaultService(), {
      shkeptncontext: 'myContext',
      eventselector: 'myId3',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.LOADED,
      artifact: 'myImage:0.0.1',
      deploymentName: 'myImage:0.0.1',
      serviceKeptnContext: 'myContext',
      evaluations: getDefaultTraces().slice(0, 1),
    });
  });

  it("should throw trace error if evaluations can't be fetched", async () => {
    mockComponentData(null, getDefaultRootTrace(), getDefaultService(), {
      shkeptncontext: 'myContext',
      eventselector: 'myId3',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.ERROR,
      kind: 'trace',
      keptnContext: 'myContext',
    });
  });

  it('should throw trace error if trace has invalid project', async () => {
    const trace = getDefaultTraces()[0];
    trace.data.project = undefined;
    mockComponentData([trace], getDefaultRootTrace(), getDefaultService(), {
      shkeptncontext: 'myContext',
      eventselector: 'myId3',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.ERROR,
      kind: 'trace',
      keptnContext: 'myContext',
    });
  });

  it('should throw trace error if trace has invalid stage', async () => {
    const trace = getDefaultTraces()[0];
    trace.data.stage = undefined;
    mockComponentData([trace], getDefaultRootTrace(), getDefaultService(), {
      shkeptncontext: 'myContext',
      eventselector: 'myId3',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.ERROR,
      kind: 'trace',
      keptnContext: 'myContext',
    });
  });

  it('should throw trace error if trace has invalid service', async () => {
    const trace = getDefaultTraces()[0];
    trace.data.service = undefined;
    mockComponentData([trace], getDefaultRootTrace(), getDefaultService(), {
      shkeptncontext: 'myContext',
      eventselector: 'myId3',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.ERROR,
      kind: 'trace',
      keptnContext: 'myContext',
    });
  });

  it('should throw default error if root trace request fails', async () => {
    mockComponentData(getDefaultTraces(), null, getDefaultService(), {
      shkeptncontext: 'myContext',
      eventselector: 'myId3',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.ERROR,
      kind: 'default',
    });
  });

  it('should throw default error if service request fails', async () => {
    mockComponentData(getDefaultTraces(), getDefaultRootTrace(), null, {
      shkeptncontext: 'myContext',
      eventselector: 'myId3',
    });
    const state = await firstValueFrom(component.state$.pipe(skip(1)));
    expect(state).toEqual({
      state: EvaluationBoardStatus.ERROR,
      kind: 'default',
    });
  });

  it('should navigate back', () => {
    createComponent();
    const locationSpy = jest.spyOn(location, 'back');
    component.goBack();
    expect(locationSpy).toHaveBeenCalled();
  });

  function mockComponentData(
    evaluations: Trace[] | null,
    rootTrace: Trace | null,
    service: Service | null,
    params: { shkeptncontext: string; eventselector?: string }
  ): void {
    jest
      .spyOn(dataService, 'getTracesByContext')
      .mockReturnValueOnce(evaluations === null ? throwError(() => new Error()) : of(evaluations));
    jest
      .spyOn(dataService, 'getTracesByContext')
      .mockReturnValueOnce(rootTrace === null ? throwError(() => new Error()) : of([rootTrace]));
    jest
      .spyOn(dataService, 'getService')
      .mockReturnValue(service === null ? throwError(() => new Error()) : of(service));
    paramMap.next(convertToParamMap(params));
    createComponent();
  }

  function getDefaultTraces(): Trace[] {
    const traces = [
      {
        id: 'myId3',
        shkeptncontext: 'myContext',
        time: '2021-10-29T08:43:15.702Z',
        type: EventTypes.EVALUATION_FINISHED,
        data: {
          stage: 'staging',
          project: 'myProject',
          service: 'myService',
        },
      },
      {
        id: 'myId2',
        shkeptncontext: 'myContext',
        time: '2021-10-29T08:43:11.702Z',
        type: EventTypes.EVALUATION_FINISHED,
        data: {
          stage: 'dev',
          project: 'myProject',
          service: 'myService',
        },
      },
      {
        id: 'myId1',
        shkeptncontext: 'myContext',
        time: '2021-10-29T08:43:10.702Z',
        type: EventTypes.EVALUATION_FINISHED,
        data: {
          stage: 'dev',
          project: 'myProject',
          service: 'myService',
        },
      },
    ];

    return traces.map((t) => Trace.fromJSON(t));
  }

  function getDefaultRootTrace(): Trace {
    return Trace.fromJSON({
      id: 'myId',
      shkeptncontext: 'myContext',
      time: '1654078902491',
      type: EventTypes.EVALUATION_FINISHED,
      data: {
        configurationChange: {
          values: {
            image: 'myRepo/myImage:0.0.1',
          },
        },
        stage: 'dev',
        project: 'myProject',
        service: 'myService',
      },
    });
  }

  function getDefaultService(): Service {
    const iService: IService = {
      serviceName: 'carts',
      creationDate: '0123456789',
      lastEventTypes: {
        [EventTypes.DEPLOYMENT_FINISHED]: {
          eventId: 'myId',
          time: '1654078902491',
          keptnContext: 'myContext',
        },
      },
    };
    return Service.fromJSON(iService);
  }

  function createComponent(): void {
    component = new EvaluationBoardComponent(
      location,
      { paramMap: paramMap.asObservable() } as ActivatedRoute,
      dataService
    );
  }
});
