import { TestBed } from '@angular/core/testing';

import { TraceDeepLinkGuard } from './trace-deep-link.guard';
import { ActivatedRouteSnapshot, convertToParamMap, ParamMap, UrlTree } from '@angular/router';
import { DataService } from '../_services/data.service';
import { firstValueFrom, Observable, of } from 'rxjs';
import { EvaluationTracesMock } from '../_services/_mockData/trace.mock';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { Trace } from '../_models/trace';

describe('TraceDeepLinkGuard', () => {
  let guard: TraceDeepLinkGuard;
  let dataService: DataService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule, RouterTestingModule],
    });
    guard = TestBed.inject(TraceDeepLinkGuard);
    dataService = TestBed.inject(DataService);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });

  it('should redirect to the latest event of a specific type', async () => {
    const keptnContext = 'fea3dc8c-5a85-435a-a86d-cee0b62f248e';
    const eventId = '214ca172-4080-4165-9a4d-39f399b17a45';

    jest.spyOn(dataService, 'getTracesByContext').mockReturnValue(of(EvaluationTracesMock));
    const url = await activateWithParams(keptnContext, EventTypes.EVALUATION_FINISHED);
    expect(url.toString()).toBe(`/project/sockshop/sequence/${keptnContext}/event/${eventId}`);
  });

  it('should redirect to sequence/:keptnContext/stage/:stage', async () => {
    const keptnContext = '218ddbfa-ed09-4cf9-887a-167a334a76d0';
    const stage = 'staging';

    jest.spyOn(dataService, 'getTracesByContext').mockReturnValue(of(EvaluationTracesMock));
    const url = await activateWithParams(keptnContext, stage);
    expect(url.toString()).toBe(`/project/sockshop/sequence/${keptnContext}/stage/${stage}`);
  });

  it('should redirect to sequence/:keptnContext', async () => {
    const keptnContext = '218ddbfa-ed09-4cf9-887a-167a334a76d0';

    jest.spyOn(dataService, 'getTracesByContext').mockReturnValue(of(EvaluationTracesMock));
    const url = await activateWithParams(keptnContext);
    expect(url.toString()).toBe(`/project/sockshop/sequence/${keptnContext}`);
  });

  it('should redirect to error if trace is invalid', async () => {
    const keptnContext = '218ddbfa-ed09-4cf9-887a-167a334a76d0';

    jest.spyOn(dataService, 'getTracesByContext').mockReturnValue(of([{} as Trace]));
    const url = await activateWithParams(keptnContext);
    expect(url.toString()).toBe(`/error?status=1001&keptnContext=${keptnContext}`);
  });

  it('should redirect to error if event selector was not found', async () => {
    const keptnContext = 'fea3dc8c-5a85-435a-a86d-cee0b62f248e';

    jest.spyOn(dataService, 'getTracesByContext').mockReturnValue(of(EvaluationTracesMock));
    const url = await activateWithParams(keptnContext, EventTypes.DEPLOYMENT_FINISHED);
    expect(url.toString()).toBe(`/error?status=1001&keptnContext=${keptnContext}`);
  });

  it('should redirect to error if keptnContext was not provided', () => {
    const url = activateWithoutKeptnContext();
    expect(url.toString()).toBe('/error?status=1001');
  });

  it('should redirect to error if no traces were found', async () => {
    const keptnContext = 'fea3dc8c-5a85-435a-a86d-cee0b62f248e';

    jest.spyOn(dataService, 'getTracesByContext').mockReturnValue(of([]));
    const url = await activateWithParams(keptnContext, EventTypes.DEPLOYMENT_FINISHED);
    expect(url.toString()).toBe(`/error?status=1001&keptnContext=${keptnContext}`);
  });

  function activateWithoutKeptnContext(): UrlTree {
    return guard.canActivate({
      get paramMap(): ParamMap {
        return convertToParamMap({});
      },
    } as ActivatedRouteSnapshot) as UrlTree;
  }

  function activateWithParams(keptnContext?: string, eventSelector?: string | EventTypes): Promise<UrlTree> {
    const response = guard.canActivate({
      get paramMap(): ParamMap {
        return convertToParamMap({
          ...(keptnContext && { keptnContext }),
          ...(eventSelector && { eventSelector }),
        });
      },
    } as ActivatedRouteSnapshot) as Observable<UrlTree>;
    return firstValueFrom(response);
  }
});
