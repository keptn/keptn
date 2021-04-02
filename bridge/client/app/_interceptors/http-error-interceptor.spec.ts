import { TestBed } from '@angular/core/testing';
import {AppModule} from "../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('HttpErrorInterceptorService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    declarations: [],
    imports: [
      AppModule,
      HttpClientTestingModule,
    ]
  }));
});
