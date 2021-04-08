import {async, fakeAsync, TestBed} from "@angular/core/testing";
import {AppModule} from "../app.module";
import {HttpClientTestingModule, HttpTestingController} from "@angular/common/http/testing";
import {ApiService} from "../_services/api.service";
import {KtbTaskItemComponent} from "../_components/ktb-task-item/ktb-task-item.component";

describe(`HttpDefaultInterceptor`, () => {
  let service: ApiService;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: []
    }).compileComponents().then(() => {
      service = TestBed.get(ApiService);
      httpMock = TestBed.get(HttpTestingController);
    });
  });

  it('should add Content-Type header', () => {
    service.getProjects().subscribe(response => {
      expect(response).toBeTruthy();
    });
    const httpRequest = httpMock.expectOne(`${service.baseUrl}/controlPlane/v1/project?disableUpstreamSync=true`);
    expect(httpRequest.request.headers.has('Content-Type')).toEqual(true);
    expect(httpRequest.request.headers.get('Content-Type')).toEqual('application/json');
    expect(httpRequest.request.headers.get('Content-Type')).toEqual('application/json');
    expect(httpRequest.request.headers.get('Content-Type')).toEqual('application/json');
    expect(httpRequest.request.headers.get('Content-Type')).toEqual('application/json');
  });
});
