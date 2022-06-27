import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbStageOverviewComponent } from './ktb-stage-overview.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbStageOverviewModule } from './ktb-stage-overview.module';
import { ProjectMock } from '../../_models/project.mock';

describe('KtbStageOverviewComponent', () => {
  let component: KtbStageOverviewComponent;
  let fixture: ComponentFixture<KtbStageOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbStageOverviewModule, HttpClientTestingModule, RouterTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbStageOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should return a link to the stage', () => {
    // given
    const project = ProjectMock;
    const stage = project.stages[0];

    // when
    const link = component.linkToStage(project, stage);

    // then
    expect(link).toEqual(['/project', 'sockshop', 'environment', 'stage', 'development']);
  });
});
