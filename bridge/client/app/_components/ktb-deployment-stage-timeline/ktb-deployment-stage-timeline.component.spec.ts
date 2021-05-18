import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbDeploymentStageTimelineComponent } from './ktb-deployment-stage-timeline.component';

describe('KtbDeploymentTimelineComponent', () => {
  let component: KtbDeploymentStageTimelineComponent;
  let fixture: ComponentFixture<KtbDeploymentStageTimelineComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbDeploymentStageTimelineComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbDeploymentStageTimelineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
