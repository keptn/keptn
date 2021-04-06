import {TestBed} from '@angular/core/testing';

import {ApiService} from './api.service';
import {HttpClientTestingModule} from "@angular/common/http/testing"
import {AppModule} from "../app.module";

describe('ApiService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    declarations: [],
    imports: [
      AppModule,
      HttpClientTestingModule,
    ],
  }));
});
