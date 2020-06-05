import { TestBed } from '@angular/core/testing';

import { HttpErrorInterceptor } from './http-error-interceptor';

describe('HttpErrorInterceptorService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should create an instance', () => {
    const service: HttpErrorInterceptor = TestBed.get(HttpErrorInterceptor);
    expect(service).toBeTruthy();
  });
});
