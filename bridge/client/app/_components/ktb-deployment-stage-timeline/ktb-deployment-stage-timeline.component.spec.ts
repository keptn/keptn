import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbDeploymentStageTimelineComponent } from './ktb-deployment-stage-timeline.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbDeploymentTimelineComponent', () => {
  let component: KtbDeploymentStageTimelineComponent;
  let fixture: ComponentFixture<KtbDeploymentStageTimelineComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbDeploymentStageTimelineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
