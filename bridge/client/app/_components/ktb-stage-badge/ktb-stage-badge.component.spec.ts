import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbStageBadgeComponent } from './ktb-stage-badge.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbStageBadgeComponent;
  let fixture: ComponentFixture<KtbStageBadgeComponent>;

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
    fixture = TestBed.createComponent(KtbStageBadgeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
