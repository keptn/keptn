import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbSubscriptionItemComponent } from './ktb-subscription-item.component';
import {HttpClientTestingModule} from '@angular/common/http/testing';
import {AppModule} from '../../app.module';
import {KtbStageDetailsComponent} from "../ktb-stage-details/ktb-stage-details.component";

describe('KtbSubscriptionItemComponent', () => {
  let component: KtbSubscriptionItemComponent;
  let fixture: ComponentFixture<KtbSubscriptionItemComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbSubscriptionItemComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
