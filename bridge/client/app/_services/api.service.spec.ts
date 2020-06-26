import { TestBed } from '@angular/core/testing';

import { ApiService } from './api.service';
import {HttpClientTestingModule} from "@angular/common/http/testing"

describe('ApiService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    declarations: [
    ],
    imports: [
      HttpClientTestingModule,
    ],
  }));

  it('should create an instance', () => {
    const service: ApiService = TestBed.get(ApiService);
    expect(service).toBeTruthy();
  });
});
