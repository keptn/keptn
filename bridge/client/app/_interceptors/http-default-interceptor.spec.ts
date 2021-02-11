import {TestBed} from "@angular/core/testing";
import {AppModule} from "../app.module";
import {HttpClientTestingModule, HttpTestingController} from "@angular/common/http/testing";
import {ApiService} from "../_services/api.service";

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
    const httpRequest = httpMock.expectOne(`${service.baseUrl}/shipyard-controller/v1/project?disableUpstreamSync=true`);
    expect(httpRequest.request.headers.has('Content-Type')).toEqual(true);
    expect(httpRequest.request.headers.get('Content-Type')).toEqual('application/json');
  });
});
