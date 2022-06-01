import { TestBed } from '@angular/core/testing';
import { AppModule } from '../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { FeatureFlagsService } from './feature-flags.service';
import { ApiService } from './api.service';
import { ApiServiceMock } from './api.service.mock';
import { DataService } from './data.service';
import { firstValueFrom } from 'rxjs';

describe('FeatureFlagsService', () => {
  let service: FeatureFlagsService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    });

    service = TestBed.inject(FeatureFlagsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should set feature flags', async () => {
    const dataService = TestBed.inject(DataService);
    dataService.loadKeptnInfo();
    const flags = await firstValueFrom(service.featureFlags$);
    expect(flags).toEqual({
      RESOURCE_SERVICE_ENABLED: false,
      D3_HEATMAP_ENABLED: false,
    });
  });
});
