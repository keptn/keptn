import { TestBed } from '@angular/core/testing';

import { OverlayService } from './overlay.service';
import { AppModule } from '../../app.module';

describe('OverlayService', () => {
  let service: OverlayService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule],
    });
    service = TestBed.inject(OverlayService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
