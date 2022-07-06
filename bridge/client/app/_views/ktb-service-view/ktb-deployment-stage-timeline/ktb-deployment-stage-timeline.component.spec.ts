import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbDeploymentStageTimelineComponent } from './ktb-deployment-stage-timeline.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbServiceViewModule } from '../ktb-service-view.module';

describe('KtbDeploymentTimelineComponent', () => {
  let component: KtbDeploymentStageTimelineComponent;
  let fixture: ComponentFixture<KtbDeploymentStageTimelineComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbServiceViewModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbDeploymentStageTimelineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
