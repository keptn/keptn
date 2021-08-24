import { TestBed } from '@angular/core/testing';
import { DataService } from './data.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../app.module';

describe('DataService', () => {
  let dataService: DataService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    });
    dataService = TestBed.inject(DataService);
  });

  it('should be an instance', () => {
    expect(dataService).toBeTruthy();
  });
});
