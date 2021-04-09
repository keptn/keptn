import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbStageBadgeComponent } from './ktb-stage-badge.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbSliBreakdownComponent} from "../ktb-sli-breakdown/ktb-sli-breakdown.component";

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbStageBadgeComponent;
  let fixture: ComponentFixture<KtbStageBadgeComponent>;

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
        fixture = TestBed.createComponent(KtbStageBadgeComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
