import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEvaluationInfoComponent } from './ktb-evaluation-info.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbEvaluationInfoComponent;
  let fixture: ComponentFixture<KtbEvaluationInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEvaluationInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
