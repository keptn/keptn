import { TestBed } from '@angular/core/testing';
import { provideMockActions } from '@ngrx/effects/testing';
import { Observable } from 'rxjs';

import { RootStoreEffects } from './root.store.effects';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { provideMockStore } from '@ngrx/store/testing';

describe(RootStoreEffects.name, () => {
  let actions$: Observable<never>;
  let effects: RootStoreEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [RootStoreEffects, provideMockActions(() => actions$), provideMockStore({})],
    });

    effects = TestBed.inject(RootStoreEffects);
  });

  it('should be created', () => {
    expect(effects).toBeTruthy();
  });
});
