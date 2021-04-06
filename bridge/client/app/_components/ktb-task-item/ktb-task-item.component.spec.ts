import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import {KtbTaskItemComponent} from './ktb-task-item.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbSubscriptionComponent} from "../ktb-subscription/ktb-subscription.component";

describe('KtbEventItemComponent', () => {
  let component: KtbTaskItemComponent;
  let fixture: ComponentFixture<KtbTaskItemComponent>;

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
        fixture = TestBed.createComponent(KtbTaskItemComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
