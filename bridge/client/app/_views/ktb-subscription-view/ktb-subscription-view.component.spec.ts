import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSubscriptionViewComponent } from './ktb-subscription-view.component';
import {AppModule} from '../../app.module';
import {HttpClient} from '@angular/common/http';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbSubscriptionViewComponent', () => {
  let component: KtbSubscriptionViewComponent;
  let fixture: ComponentFixture<KtbSubscriptionViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSubscriptionViewComponent ],
      imports: [AppModule, HttpClientTestingModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSubscriptionViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
