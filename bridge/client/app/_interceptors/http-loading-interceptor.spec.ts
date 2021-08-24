import { TestBed } from '@angular/core/testing';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { Overlay } from '@angular/cdk/overlay';
import { HttpLoadingInterceptor } from './http-loading-interceptor';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('HttpErrorInterceptorService', () => {

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        HttpClientTestingModule,
      ],
      providers: [
        Overlay,
        {
          provide: HTTP_INTERCEPTORS,
          useClass: HttpLoadingInterceptor,
          multi: true,
        },
      ],
    }).compileComponents();
  });

  it('should be an instance', () => {
    const interceptor = TestBed.inject(HttpLoadingInterceptor);
    expect(interceptor).toBeTruthy();
  });
});
