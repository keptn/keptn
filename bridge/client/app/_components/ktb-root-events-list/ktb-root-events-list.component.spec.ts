import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbRootEventsListComponent } from './ktb-root-events-list.component';
import { KtbEventsListComponent } from '../ktb-events-list/ktb-events-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';
import { Project } from '../../_models/project';
import { By } from '@angular/platform-browser';

describe('KtbEventsListComponent', () => {
  let component: KtbRootEventsListComponent;
  let fixture: ComponentFixture<KtbRootEventsListComponent>;
  let dataService: DataService;
  const projectName = 'sockshop';
  let project: Project;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        {
          provide: DataService,
          useClass: DataServiceMock,
        },
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({}),
            params: of({ projectName }),
            queryParams: of({}),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbRootEventsListComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    dataService.loadProjects(); // reset project.sequences
    // @ts-ignore
    dataService.getProject(projectName).subscribe((pr: Project) => {
      project = pr;
      fixture.detectChanges();
    });
  });

  it('should create root-events-list component', () => {
    expect(component).toBeTruthy();
  });

  it('should show 25 sequences', () => {
    // given
    dataService.loadSequences(project);
    component.events = project.sequences;
    fixture.detectChanges();

    // then
    const sequences = fixture.nativeElement.querySelectorAll('ktb-selectable-tile');
    const showMoreButton = fixture.nativeElement.querySelector('button[dt-show-more]');
    expect(sequences.length).toEqual(25);
    expect(showMoreButton).toBeTruthy();
    expect(component.events).toEqual(project.sequences);
  });

  it('should load old sequences', () => {
    // given
    dataService.loadSequences(project);
    component.events = project.sequences;
    fixture.detectChanges();

    // when
    component.loadOldSequences();
    component.events = project.sequences;
    fixture.detectChanges();

    // then
    const sequences = fixture.nativeElement.querySelectorAll('ktb-selectable-tile');
    const showMoreButton = fixture.nativeElement.querySelector('button[dt-show-more]');
    expect(sequences.length).toEqual(35);
    expect(showMoreButton).toBeTruthy();
  });

  it('should not show "show older sequences"', () => {
    // given
    dataService.loadSequences(project);
    component.events = project.sequences;
    component.loadOldSequences();
    component.events = project.sequences;

    // when
    component.loadOldSequences();
    component.events = project.sequences;
    fixture.detectChanges();

    // then
    const sequences = fixture.nativeElement.querySelectorAll('ktb-selectable-tile');
    const showMoreButton = fixture.nativeElement.querySelector('button[dt-show-more]');
    expect(sequences.length).toEqual(36);
    expect(showMoreButton).toBeFalsy();
  });

  it('should select provided sequence', () => {
    // given
    const selectedSequenceIndex = 1;
    dataService.loadSequences(project);
    component.events = project.sequences;
    component.selectedEvent = project.sequences[selectedSequenceIndex];
    fixture.detectChanges();

    // then
    const selectedTile = getSequenceTile(selectedSequenceIndex);
    expect(selectedTile.getAttribute('class')).toContain('ktb-tile-selected');
  });

  it('should select sequence', () => {
    // given
    const selectedSequenceIndex = 5;
    const changeEvent = jest.spyOn(component.selectedEventChange, 'emit');
    dataService.loadSequences(project);
    component.events = project.sequences;
    fixture.detectChanges();

    // when
    const targetSequence = getSequenceTile(selectedSequenceIndex);
    const eventData = { sequence: project.sequences[selectedSequenceIndex], stage: undefined };
    targetSequence.click();
    fixture.detectChanges();

    // then
    expect(targetSequence.getAttribute('class')).toContain('ktb-tile-selected');
    expect(component.selectedEvent).toEqual(eventData.sequence);
    expect(changeEvent).toHaveBeenCalledWith(eventData);
  });

  it('should select sequence with stage', () => {
    // given
    const selectedSequenceIndex = 8;
    const changeEvent = jest.spyOn(component.selectedEventChange, 'emit');
    dataService.loadSequences(project);
    component.events = project.sequences;
    fixture.detectChanges();

    // when
    const targetSequence = getSequenceTile(selectedSequenceIndex);
    const stageBadges = targetSequence.querySelectorAll('ktb-stage-badge');
    const targetStage = stageBadges[0];
    const stageName = targetStage.querySelector('dt-tag').textContent;
    targetStage.click();
    fixture.detectChanges();

    // then
    expect(stageBadges.length).toEqual(2);
    expect(targetSequence.getAttribute('class')).toContain('ktb-tile-selected');
    expect(changeEvent).toHaveBeenCalledWith({ sequence: project.sequences[selectedSequenceIndex], stage: stageName });
  });

  it('should have a no specific class when a sequence is running', () => {
    // given
    prepareSequenceElement(false, false, false);

    // when
    const sequence = fixture.debugElement.queryAll(By.css('.ktb-selectable-tile'))[0];
    sequence.nativeElement.click();
    fixture.detectChanges();

    // then
    expect(sequence.classes['ktb-tile-selected']).toBe(true);
  });

  it('should have an error class when a sequence is finished and failed', () => {
    // given
    prepareSequenceElement(true, true, false);

    // when
    const sequence = fixture.debugElement.queryAll(By.css('.ktb-selectable-tile'))[0];
    fixture.detectChanges();

    // then
    expect(sequence.classes['ktb-tile-error']).toBe(true);
  });

  it('should have a success class when a sequence is finished and not failed', () => {
    // given
    prepareSequenceElement(true, false, false);

    // when
    const sequence = fixture.debugElement.queryAll(By.css('.ktb-selectable-tile'))[0];
    fixture.detectChanges();

    // then
    expect(sequence.classes['ktb-tile-success']).toBe(true);
  });

  it('should have a highlight class when a sequence has pending approvals', () => {
    // given
    prepareSequenceElement(false, false, true);

    // when
    const sequence = fixture.debugElement.queryAll(By.css('.ktb-selectable-tile'))[0];
    fixture.detectChanges();

    // then
    expect(sequence.classes['ktb-tile-highlight']).toBe(true);
  });


  function getSequenceTile(index: number): HTMLElement {
    return fixture.nativeElement.querySelector(
      `ktb-selectable-tile[uitestid="keptn-root-events-list-${project.sequences[index].shkeptncontext}"]`
    );
  }

  function prepareSequenceElement(isFinished: boolean, isFaulty: boolean, hasPendingApproval: boolean): void {
    dataService.loadSequences(project);
    component.events = project.sequences;
    jest.spyOn(component.events[0], 'isFinished').mockReturnValue(isFinished);
    jest.spyOn(component.events[0], 'isFaulty').mockReturnValue(isFaulty);
    jest.spyOn(component.events[0], 'hasPendingApproval').mockReturnValue(hasPendingApproval);
    fixture.detectChanges();
  }
});
