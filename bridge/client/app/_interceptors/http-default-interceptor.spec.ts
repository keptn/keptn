import { TestBed } from '@angular/core/testing';
import { AppModule } from '../app.module';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { ApiService } from '../_services/api.service';

describe(`HttpDefaultInterceptor`, () => {
  let service: ApiService;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    service = TestBed.inject(ApiService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should add Content-Type header', () => {
    service.getProjects().subscribe((response) => {
      expect(response).toBeTruthy();
    });
    const httpRequest = httpMock.expectOne(`${service.baseUrl}/controlPlane/v1/project?disableUpstreamSync=true`);
    expect(httpRequest.request.headers.has('Content-Type')).toEqual(true);
    expect(httpRequest.request.headers.get('Content-Type')).toEqual('application/json');
  });
});
