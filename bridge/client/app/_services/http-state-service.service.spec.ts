import { TestBed } from '@angular/core/testing';

import { HttpStateService } from './http-state.service';

describe('HttpStateServiceService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: HttpStateService = TestBed.get(HttpStateService);
    expect(service).toBeTruthy();
  });
});
