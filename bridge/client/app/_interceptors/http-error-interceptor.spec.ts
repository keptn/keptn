import { TestBed } from '@angular/core/testing';
import { HTTP_INTERCEPTORS, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { HttpErrorInterceptor } from './http-error-interceptor';
import { Overlay } from '@angular/cdk/overlay';
import { HttpClientTestingModule, HttpTestingController, TestRequest } from '@angular/common/http/testing';
import { RETRY_ON_HTTP_ERROR } from '../_utils/app.utils';
import { ApiService } from '../_services/api.service';
import { NotificationsService } from '../_services/notifications.service';
import { NotificationType } from '../_models/notification';
import { SecretScopeDefault } from '../../../shared/interfaces/secret-scope';
import { Secret } from '../_models/secret';

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
        { provide: RETRY_ON_HTTP_ERROR, useValue: false },
      ],
    });

    httpErrorInterceptor = TestBed.inject(HttpErrorInterceptor);
    httpMock = TestBed.inject(HttpTestingController);
    apiService = TestBed.inject(ApiService);
  });

  it('should be an instance', () => {
    expect(httpErrorInterceptor).toBeTruthy();
  });

  it('should show error for HttpErrorResponse', () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    apiService.getMetadata().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: HttpErrorResponse = new HttpErrorResponse({ status: 500 });
    testRequest.flush('server error', errorEvent);

    // then
    expect(spy).toHaveBeenCalledWith(NotificationType.ERROR, 'server error');
  });

  it('should show an error for client-side error', () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    apiService.getMetadata().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: ErrorEvent = new ErrorEvent('client error');
    testRequest.error(errorEvent, { status: 404 });

    // then
    expect(spy).toHaveBeenCalledWith(NotificationType.ERROR, 'Http failure response for ./api/v1/metadata: 404 ');
  });

  it('should show a generic error notification when unauthorized', () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    apiService.getMetadata().subscribe();
    apiService.getProjects().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: HttpErrorResponse = new HttpErrorResponse({ status: 401 });
    testRequest.flush('', errorEvent);

    const testRequestProjects: TestRequest = httpMock.expectOne(
      './api/controlPlane/v1/project?disableUpstreamSync=true'
    );
    testRequestProjects.flush('', errorEvent);

    // then
    expect(spy).toBeCalledTimes(1);
    expect(spy).toHaveBeenCalledWith(NotificationType.ERROR, 'Could not authorize.');
  });

  it('should show error notification when incorrect api key', () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    apiService.getMetadata().subscribe();
    apiService.getProjects().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: HttpErrorResponse = new HttpErrorResponse({ status: 401 });
    testRequest.flush('incorrect api key auth', errorEvent);

    const testRequestProjects: TestRequest = httpMock.expectOne(
      './api/controlPlane/v1/project?disableUpstreamSync=true'
    );
    testRequestProjects.flush('incorrect api key auth', errorEvent);

    // then
    expect(spy).toBeCalledTimes(1);
    expect(spy).toHaveBeenCalledWith(
      NotificationType.ERROR,
      'Could not authorize API token. Please check the configured API token.'
    );
  });

  it('should show error notification when basic auth is unauthorized', () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    apiService.getMetadata().subscribe();
    apiService.getProjects().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: ErrorEvent = new ErrorEvent('', { error: {} });
    const headers = new HttpHeaders({ 'keptn-auth-type': 'BASIC' });
    testRequest.error(errorEvent, { headers, status: 401 });

    const testRequestProjects: TestRequest = httpMock.expectOne(
      './api/controlPlane/v1/project?disableUpstreamSync=true'
    );
    testRequestProjects.error(errorEvent, { headers, status: 401 });

    // then
    expect(spy).toBeCalledTimes(1);
    expect(spy).toHaveBeenCalledWith(
      NotificationType.ERROR,
      'Login credentials invalid. Please check your provided username and password.'
    );
  });

  it('should show info notification when oauth is unauthorized', () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    apiService.getMetadata().subscribe();
    apiService.getProjects().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: ErrorEvent = new ErrorEvent('', { error: {} });
    const headers = new HttpHeaders({ 'keptn-auth-type': 'OAUTH' });
    testRequest.error(errorEvent, { headers, status: 401 });

    const testRequestProjects: TestRequest = httpMock.expectOne(
      './api/controlPlane/v1/project?disableUpstreamSync=true'
    );
    testRequestProjects.error(errorEvent, { headers, status: 401 });

    // then
    expect(spy).toBeCalledTimes(1);
    expect(spy).toHaveBeenCalledWith(NotificationType.INFO, 'Login required. Redirecting to login.');
  });

  it('should show error notification in case of 403', () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    apiService.getMetadata().subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/v1/metadata');
    const errorEvent: ErrorEvent = new ErrorEvent('', { error: {} });
    testRequest.error(errorEvent, { status: 403 });

    // then
    expect(spy).toHaveBeenCalledWith(NotificationType.ERROR, 'You do not have the permissions to perform this action.');
  });

  it('should not show any notification when a secret already exists', () => {
    // given
    const spy = jest.spyOn(TestBed.inject(NotificationsService), 'addNotification');

    const secret = new Secret();
    secret.name = 'secret';
    secret.scope = SecretScopeDefault.DEFAULT;

    apiService.addSecret(secret).subscribe();

    const testRequest: TestRequest = httpMock.expectOne('./api/secrets/v1/secret');
    const errorEvent: ErrorEvent = new ErrorEvent('', { error: {} });
    testRequest.error(errorEvent, { status: 409 });

    // when

    // then
    expect(spy).toHaveBeenCalledTimes(0);
  });
});
