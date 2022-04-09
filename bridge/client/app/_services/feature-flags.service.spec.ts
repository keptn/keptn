import { TestBed } from '@angular/core/testing';
import { AppModule } from '../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { FeatureFlagsService } from './feature-flags.service';
import { APIService } from './api.service';
import { ApiServiceMock } from './api.service.mock';
import { DataService } from './data.service';

describe('FeatureFlagsService', () => {
  let service: FeatureFlagsService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: APIService, useClass: ApiServiceMock }],
    });

    service = TestBed.inject(FeatureFlagsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should set feature flags', () => {
    const dataService = TestBed.inject(DataService);
    dataService.loadKeptnInfo();
    expect(service.featureFlags).toEqual({});
  });
});
