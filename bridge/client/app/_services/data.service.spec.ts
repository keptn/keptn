import { TestBed } from '@angular/core/testing';
import { DataService } from './data.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiService } from './api.service';
import { TriggerSequenceData } from '../_models/trigger-sequence';
import moment from 'moment';
import { firstValueFrom, of, throwError } from 'rxjs';
import { Service } from '../_models/service';
import { IService } from '../../../shared/interfaces/service';
import { Trace } from '../_models/trace';
import { HttpErrorResponse, HttpHeaders, HttpResponse } from '@angular/common/http';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { ApiServiceMock } from './api.service.mock';
import { filter } from 'rxjs/operators';

describe('DataService', () => {
  let dataService: DataService;
  let apiService: ApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    });
    dataService = TestBed.inject(DataService);
    apiService = TestBed.inject(ApiService);
  });

  it('should be an instance', () => {
    expect(dataService).toBeTruthy();
  });

  it('should load projects', (done) => {
    dataService.loadProjects();
    dataService.projects.pipe(filter((projects) => !!projects && projects?.length > 0)).subscribe((projects) => {
      expect(projects?.length).toBe(3);
      done();
    });
  });

  it('should trigger a delivery', () => {
    // given
    const spy = jest.spyOn(apiService, 'triggerSequence');
    const data: TriggerSequenceData = {
      project: 'podtato-head',
      stage: 'hardening',
      service: 'helloservice',
      configurationChange: {
        values: {
          image: 'docker.io/keptn:v0.1.2',
        },
      },
    };

    // when
    dataService.triggerDelivery(data);

    // then
    expect(spy).toHaveBeenCalledWith('sh.keptn.event.hardening.delivery.triggered', data);
  });

  it('should trigger an evaluation', () => {
    // given
    const spy = jest.spyOn(apiService, 'triggerSequence');
    const date = moment().toISOString();
    const data: TriggerSequenceData = {
      project: 'podtato-head',
      stage: 'hardening',
      service: 'helloservice',
      evaluation: {
        timeframe: '1h15m',
        start: date,
      },
    };

    // when
    dataService.triggerEvaluation(data);

    // then
    expect(spy).toHaveBeenCalledWith('sh.keptn.event.hardening.evaluation.triggered', data);
  });

  it('should trigger a custom sequence', () => {
    // given
    const spy = jest.spyOn(apiService, 'triggerSequence');
    const data: TriggerSequenceData = {
      project: 'podtato-head',
      stage: 'hardening',
      service: 'helloservice',
      labels: {
        key1: 'val1',
      },
    };

    // when
    dataService.triggerCustomSequence(data, 'testsequence');

    // then
    expect(spy).toHaveBeenCalledWith('sh.keptn.event.hardening.testsequence.triggered', data);
  });

  it('should map response to service', async () => {
    const iService: IService = {
      serviceName: 'carts',
      creationDate: '0123456789',
      lastEventTypes: {},
    };
    jest.spyOn(apiService, 'getService').mockReturnValue(of(iService));
    const service = await firstValueFrom(dataService.getService('sockshop', 'dev', 'carts'));
    expect(service).toBeInstanceOf(Service);
  });

  it('should get traces by context', async () => {
    setGetTracesResponse([getDefaultTrace() as Trace]);
    const traces = await firstValueFrom(dataService.getTracesByContext('abc'));
    for (const trace of traces) {
      expect(trace).toBeInstanceOf(Trace);
    }
  });

  it('should send an approval once', () => {
    const sendApprovalSpy = jest.spyOn(apiService, 'sendApprovalEvent');
    sendApprovalSpy.mockReturnValue(of({}));
    dataService.sendApprovalEvent(Trace.fromJSON({}), true).subscribe();
    expect(sendApprovalSpy).toHaveBeenCalledTimes(1);
  });

  it('should load uniform log status', async () => {
    // given
    jest.spyOn(apiService, 'hasUnreadUniformRegistrationLogs').mockReturnValue(of(true));

    // when
    dataService.loadUnreadUniformRegistrationLogs();

    // then
    expect(await firstValueFrom(dataService.hasUnreadUniformRegistrationLogs)).toBe(true);
  });

  it('should not load uniform log status', () => {
    // given
    const loadLogSpy = jest.spyOn(apiService, 'hasUnreadUniformRegistrationLogs');
    dataService.setHasUnreadUniformRegistrationLogs(true);

    // when
    dataService.loadUnreadUniformRegistrationLogs();

    // then
    expect(loadLogSpy).not.toHaveBeenCalled();
  });

  it('should return traces by given ids', async () => {
    const result = await firstValueFrom(dataService.getTracesByIds('myProject', ['id1', 'id2']));

    expect(result[0]).toBeInstanceOf(Trace);
  });

  it('should return default bridgeInfo if request failed', async () => {
    const errorResponse = new HttpErrorResponse({
      headers: new HttpHeaders({ 'keptn-auth-type': 'OAUTH' }),
      status: 403,
    });
    jest.spyOn(apiService, 'getKeptnInfo').mockReturnValue(throwError(() => errorResponse));
    dataService.loadKeptnInfo();
    const keptnInfo = await firstValueFrom(dataService.keptnInfo);
    expect(keptnInfo?.bridgeInfo).toEqual({
      authType: 'OAUTH',
      cliDownloadLink: '',
      enableVersionCheckFeature: false,
      featureFlags: {},
      showApiToken: false,
    });
  });

  function setGetTracesResponse(traces: Trace[]): void {
    const response = new HttpResponse({ body: { events: traces, totalCount: 0, pageSize: 1, nextPageKey: 1 } });
    jest.spyOn(apiService, 'getTraces').mockReturnValue(of(response));
  }

  function getDefaultTrace(): unknown {
    return {
      id: 'myId',
      shkeptncontext: 'myContext',
      time: '123456789',
      type: EventTypes.EVALUATION_TRIGGERED,
      data: {
        stage: 'myStage',
        project: 'myProject',
        service: 'myService',
      },
    };
  }
});
