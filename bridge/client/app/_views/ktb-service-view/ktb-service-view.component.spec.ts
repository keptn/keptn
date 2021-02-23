import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {KtbServiceViewComponent} from './ktb-service-view.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEventsListComponent', () => {
  let component: KtbServiceViewComponent;
  let fixture: ComponentFixture<KtbServiceViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbServiceViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
