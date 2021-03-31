import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSubscriptionComponent } from './ktb-subscription.component';
import {HttpClientTestingModule} from '@angular/common/http/testing';
import {AppModule} from '../../app.module';

describe('KtbSubscriptionComponent', () => {
  let component: KtbSubscriptionComponent;
  let fixture: ComponentFixture<KtbSubscriptionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSubscriptionComponent ],
      imports: [AppModule, HttpClientTestingModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSubscriptionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
