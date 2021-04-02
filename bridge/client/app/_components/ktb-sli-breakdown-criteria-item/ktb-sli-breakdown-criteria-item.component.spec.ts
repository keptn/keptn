import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSliBreakdownCriteriaItemComponent } from './ktb-sli-breakdown-criteria-item.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbSliBreakdownCriteriaItemComponent', () => {
  let component: KtbSliBreakdownCriteriaItemComponent;
  let fixture: ComponentFixture<KtbSliBreakdownCriteriaItemComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSliBreakdownCriteriaItemComponent ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSliBreakdownCriteriaItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
