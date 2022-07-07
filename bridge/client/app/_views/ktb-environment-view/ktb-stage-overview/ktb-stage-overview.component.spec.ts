import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbStageOverviewComponent } from './ktb-stage-overview.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbEnvironmentViewModule } from '../ktb-environment-view.module';

describe('KtbStageOverviewComponent', () => {
  let component: KtbStageOverviewComponent;
  let fixture: ComponentFixture<KtbStageOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbEnvironmentViewModule, HttpClientTestingModule, RouterTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbStageOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
