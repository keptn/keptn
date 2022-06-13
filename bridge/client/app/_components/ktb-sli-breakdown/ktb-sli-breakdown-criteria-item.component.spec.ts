import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSliBreakdownCriteriaItemComponent } from './ktb-sli-breakdown-criteria-item.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbSliBreakdownModule } from './ktb-sli-breakdown.module';

describe('KtbSliBreakdownCriteriaItemComponent', () => {
  let component: KtbSliBreakdownCriteriaItemComponent;
  let fixture: ComponentFixture<KtbSliBreakdownCriteriaItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSliBreakdownModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSliBreakdownCriteriaItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
