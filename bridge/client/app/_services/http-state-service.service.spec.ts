import { TestBed } from '@angular/core/testing';

import { HttpStateServiceService } from './http-state.service';

describe('HttpStateServiceService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: HttpStateServiceService = TestBed.get(HttpStateServiceService);
    expect(service).toBeTruthy();
  });
});
