import { TestBed } from '@angular/core/testing';
import { HTTP_INTERCEPTORS, HttpHeaders } from '@angular/common/http';
import { HttpErrorInterceptor } from './http-error-interceptor';
import { Overlay } from '@angular/cdk/overlay';
import { HttpClientTestingModule, HttpTestingController, TestRequest } from '@angular/common/http/testing';
import { RETRY_ON_HTTP_ERROR } from '../_utils/app.utils';
import { ApiService } from '../_services/api.service';
import { NotificationsService } from '../_services/notifications.service';
import { DtToast } from '@dynatrace/barista-components/toast';

describe('HttpErrorInterceptorService', () => {
  let httpErrorInterceptor: HttpErrorInterceptor;
  let httpMock: HttpTestingController;
  let apiService: ApiService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [
        Overlay,
        {
          provide: HTTP_INTERCEPTORS,
          useClass: HttpErrorInterceptor,
          multi: true,
        },
        {provide: RETRY_ON_HTTP_ERROR, useValue: false},
      ],
    });

    httpErrorInterceptor = TestBed.inject(HttpErrorInterceptor);
    httpMock = TestBed.inject(HttpTestingController);
    apiService = TestBed.inject(ApiService);
  });

  it('should be an instance', () => {
    expect(httpErrorInterceptor).toBeTruthy();
  });

  it('should show an error when any other error than 401 is returned', async () => {
    // given
    const spy = jest.spyOn(TestBed.inject(DtToast), 'create');

    await apiService.getMetadata().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: ErrorEvent = new ErrorEvent('', {error: {}});
    testRequest.error(errorEvent, {status: 404});

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should show an generic error when unauthorized', async () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');


    await apiService.getMetadata().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: ErrorEvent = new ErrorEvent('', {error: {}});
    testRequest.error(errorEvent, {status: 401});

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should show a toast when oauth is redirected', async () => {
    // given
    const spy = jest.spyOn(TestBed.inject(DtToast), 'create');

    await apiService.getMetadata().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: ErrorEvent = new ErrorEvent('', {error: {}});
    // @ts-ignore
    const headers = new HttpHeaders({'keptn-auth-type': 'OAUTH'});
    testRequest.error(errorEvent, {headers, status: 401});

    // then
    expect(spy).toHaveBeenCalledWith('Login required. Redirecting to login.');
  });

  it('should show a error notification when basic auth is unauthorized', async () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    await apiService.getMetadata().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: ErrorEvent = new ErrorEvent('', {error: {}});
    // @ts-ignore
    const headers = new HttpHeaders({'keptn-auth-type': 'BASIC'});
    testRequest.error(errorEvent, {headers, status: 401});

    // then
    expect(spy).toHaveBeenCalledWith('error', 'Login credentials invalid. Please check your provided username and password.');
  });
});
