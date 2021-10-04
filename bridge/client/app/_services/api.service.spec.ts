import { TestBed } from '@angular/core/testing';
import { ApiService } from './api.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../app.module';

describe('ApiService', () => {
  let apiService: ApiService;
  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
    });
    apiService = TestBed.inject(ApiService);
  });

  it('should be an instance', () => {
    expect(apiService).toBeTruthy();
  });
});
