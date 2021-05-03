import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbDeploymentStageTimelineComponent } from './ktb-deployment-stage-timeline.component';

describe('KtbDeploymentTimelineComponent', () => {
  let component: KtbDeploymentStageTimelineComponent;
  let fixture: ComponentFixture<KtbDeploymentStageTimelineComponent>;

  beforeEach(async(() => {
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
