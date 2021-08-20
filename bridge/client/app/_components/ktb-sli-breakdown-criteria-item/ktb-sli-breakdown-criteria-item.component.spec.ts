import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSliBreakdownCriteriaItemComponent } from './ktb-sli-breakdown-criteria-item.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbSliBreakdownCriteriaItemComponent', () => {
  let component: KtbSliBreakdownCriteriaItemComponent;
  let fixture: ComponentFixture<KtbSliBreakdownCriteriaItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSliBreakdownCriteriaItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
