import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSliBreakdownComponent } from './ktb-sli-breakdown.component';
import {KtbEvaluationDetailsComponent} from "../ktb-evaluation-details/ktb-evaluation-details.component";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbSliBreakdownComponent;
  let fixture: ComponentFixture<KtbSliBreakdownComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSliBreakdownComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
