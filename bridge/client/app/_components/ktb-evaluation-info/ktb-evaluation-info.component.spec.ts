import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbEvaluationInfoComponent } from './ktb-evaluation-info.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbEvaluationInfoComponent;
  let fixture: ComponentFixture<KtbEvaluationInfoComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEvaluationInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
