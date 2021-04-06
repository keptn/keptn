import {TestBed} from '@angular/core/testing';

import {HttpStateService} from './http-state.service';
import {AppModule} from "../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('HttpStateServiceService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    declarations: [],
    imports: [
      AppModule,
      HttpClientTestingModule,
    ]
  }));
});
