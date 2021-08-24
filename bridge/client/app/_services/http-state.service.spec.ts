import { TestBed } from '@angular/core/testing';
import { AppModule } from '../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { HttpStateService } from './http-state.service';

describe('HttpStateService', () => {
  let service: HttpStateService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    });

    service = TestBed.inject(HttpStateService);
  });

  it('should be an instance', () => {
    expect(service).toBeTruthy();
  });
});
