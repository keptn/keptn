import { TestBed } from '@angular/core/testing';

import { DataService } from './data.service';
import {HttpClientTestingModule} from "@angular/common/http/testing"

describe('DataService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    declarations: [
    ],
    imports: [
      HttpClientTestingModule,
    ],
  }));

  it('should create an instance', () => {
    const service: DataService = TestBed.get(DataService);
    expect(service).toBeTruthy();
  });
});
