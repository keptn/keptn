import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RETRY_ON_HTTP_ERROR } from '../../../_utils/app.utils';
import { KtbProjectListComponent } from './ktb-project-list.component';
import { ProjectsMock } from '../../../_services/_mockData/projects.mock';
import { SequencesMock } from '../../../_services/_mockData/sequences.mock';
import { KtbDashboardLegacyViewModule } from '../ktb-dashboard-legacy-view.module';

describe('KtbProjectListComponent', () => {
  let component: KtbProjectListComponent;
  let fixture: ComponentFixture<KtbProjectListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbDashboardLegacyViewModule, HttpClientTestingModule],
      providers: [{ provide: RETRY_ON_HTTP_ERROR, useValue: false }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should extract the sequences of a project', () => {
    // given
    const project1 = ProjectsMock[0];
    const project2 = ProjectsMock[1];
    const project3 = ProjectsMock[2];
    component.sequences = { [project2.projectName]: [], [project1.projectName]: SequencesMock };

    // when
    const sequences1 = component.getSequencesPerProject(project1);
    const sequences2 = component.getSequencesPerProject(project2);
    const sequences3 = component.getSequencesPerProject(project3);

    // then
    expect(sequences1.length).toEqual(36);
    expect(sequences2.length).toEqual(0);
    expect(sequences3.length).toEqual(0);
  });
});
