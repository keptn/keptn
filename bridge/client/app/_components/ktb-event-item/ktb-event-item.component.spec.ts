import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';

import {KtbEventItemComponent} from './ktb-event-item.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbEvaluationDetailsComponent} from "../ktb-evaluation-details/ktb-evaluation-details.component";

describe('KtbEventItemComponent', () => {
  let component: KtbEventItemComponent;
  let fixture: ComponentFixture<KtbEventItemComponent>;

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
        fixture = TestBed.createComponent(KtbEventItemComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
