import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbSliBreakdownComponent } from './ktb-sli-breakdown.component';
import {KtbEvaluationDetailsComponent} from "../ktb-evaluation-details/ktb-evaluation-details.component";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";
import {KtbSliBreakdownCriteriaItemComponent} from "../ktb-sli-breakdown-criteria-item/ktb-sli-breakdown-criteria-item.component";

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbSliBreakdownComponent;
  let fixture: ComponentFixture<KtbSliBreakdownComponent>;

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
        fixture = TestBed.createComponent(KtbSliBreakdownComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
