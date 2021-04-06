import { TestBed } from '@angular/core/testing';

import { NotificationsService } from './notifications.service';
import {AppModule} from "../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('NotificationsService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    declarations: [],
    imports: [
      AppModule,
      HttpClientTestingModule,
    ]
  }));

  it('should be created', () => {
    const service: NotificationsService = TestBed.get(NotificationsService);
    expect(service).toBeTruthy();
  });
});
