import { TestBed } from '@angular/core/testing';
import { NotificationsService } from './notifications.service';
import { AppModule } from '../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('NotificationsService', () => {
  let service: NotificationsService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
    });

    service = TestBed.inject(NotificationsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
