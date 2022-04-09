import { TestBed } from '@angular/core/testing';
import { APIService } from './api.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../app.module';

describe('APIService', () => {
  let apiService: APIService;
  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
    });
    apiService = TestBed.inject(APIService);
  });

  it('should be an instance', () => {
    expect(apiService).toBeTruthy();
  });
});
