import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSliBreakdownCriteriaItemComponent } from './ktb-sli-breakdown-criteria-item.component';

describe('KtbSliBreakdownCriteriaItemComponent', () => {
  let component: KtbSliBreakdownCriteriaItemComponent;
  let fixture: ComponentFixture<KtbSliBreakdownCriteriaItemComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSliBreakdownCriteriaItemComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSliBreakdownCriteriaItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
