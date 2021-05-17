import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbServiceDetailsComponent } from './ktb-service-details.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbServiceDetailsComponent', () => {
  let component: KtbServiceDetailsComponent;
  let fixture: ComponentFixture<KtbServiceDetailsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbServiceDetailsComponent ],
      imports: [AppModule, HttpClientTestingModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbServiceDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
