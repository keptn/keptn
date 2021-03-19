import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSubscriptionItemComponent } from './ktb-subscription-item.component';
import {HttpClientTestingModule} from '@angular/common/http/testing';
import {AppModule} from '../../app.module';

describe('KtbSubscriptionItemComponent', () => {
  let component: KtbSubscriptionItemComponent;
  let fixture: ComponentFixture<KtbSubscriptionItemComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSubscriptionItemComponent ],
      imports: [AppModule, HttpClientTestingModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSubscriptionItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
