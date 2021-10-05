import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbStageOverviewComponent } from './ktb-stage-overview.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbStageOverviewComponent', () => {
  let component: KtbStageOverviewComponent;
  let fixture: ComponentFixture<KtbStageOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbStageOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
