import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEvaluationDetailsComponent } from './ktb-evaluation-details.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbEvaluationDetailsComponent;
  let fixture: ComponentFixture<KtbEvaluationDetailsComponent>;

  beforeEach(async(() => {
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
    fixture = TestBed.createComponent(KtbEvaluationDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
