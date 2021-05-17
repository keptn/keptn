import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbSubscriptionComponent } from './ktb-subscription.component';
import {HttpClientTestingModule} from '@angular/common/http/testing';
import {AppModule} from '../../app.module';
import {KtbSubscriptionItemComponent} from "../ktb-subscription-item/ktb-subscription-item.component";

describe('KtbSubscriptionComponent', () => {
  let component: KtbSubscriptionComponent;
  let fixture: ComponentFixture<KtbSubscriptionComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbSubscriptionComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
