import {ComponentFixture, TestBed} from '@angular/core/testing';

import {KtbSequenceControlsComponent} from './ktb-sequence-controls.component';
import {AppModule} from '../../app.module';
import { DataService } from '../../_services/data.service';
import { Project } from '../../_models/project';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataServiceMock } from '../../_services/data.service.mock';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';

describe('KtbSequenceControlsComponent', () => {
  let component: KtbSequenceControlsComponent;
  let fixture: ComponentFixture<KtbSequenceControlsComponent>;
  let dataService: DataService;
  const projectName = 'sockshop';
  let project: Project;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {
          provide: DataService,
          useClass: DataServiceMock,
        },
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({}),
            params: of({projectName}),
            queryParams: of({}),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceControlsComponent);
    component = fixture.componentInstance;

    dataService = fixture.debugElement.injector.get(DataService);
    dataService.loadProjects(); // reset project.sequences
    // @ts-ignore
    dataService.getProject(projectName).subscribe((pr: Project) => {
      project = pr;
      fixture.detectChanges();
    });
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show pause and abort for started sequence', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[6];
    fixture.detectChanges();

    const sequencePauseButton = fixture.nativeElement.querySelector('button[uitestid=sequencePauseButton]');
    const sequenceResumeButton = fixture.nativeElement.querySelector('button[uitestid=sequenceResumeButton]');
    const sequenceAbortButton = fixture.nativeElement.querySelector('button[uitestid=sequenceAbortButton]');

    // then
    expect(component.sequence.isFinished()).toBe(false);
    expect(component.sequence.isPaused()).toBe(false);

    expect(sequencePauseButton).toBeTruthy();
    expect(sequenceResumeButton).toBeFalsy();
    expect(sequenceAbortButton).toBeTruthy();

    expect(sequencePauseButton.disabled).toBeFalsy();
    expect(sequenceAbortButton.disabled).toBeFalsy();
  });

  it('should pause sequence on click', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[6];
    fixture.detectChanges();

    const sequencePauseButton = fixture.nativeElement.querySelector('button[uitestid=sequencePauseButton]');
    const spyPause = jest.spyOn(component, 'triggerPauseSequence');

    // when
    sequencePauseButton.dispatchEvent(new Event('click'));

    // then
    expect(spyPause).toHaveBeenCalled();
  });

  it('should abort sequence on click', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[6];
    fixture.detectChanges();

    const sequenceAbortButton = fixture.nativeElement.querySelector('button[uitestid=sequenceAbortButton]');
    const spyAbort = jest.spyOn(component, 'triggerAbortSequence');

    // when
    sequenceAbortButton.dispatchEvent(new Event('click'));

    // then
    expect(spyAbort).toHaveBeenCalled();
  });

  it('should show resume and abort for paused sequence', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[3];
    fixture.detectChanges();

    const sequencePauseButton = fixture.nativeElement.querySelector('button[uitestid=sequencePauseButton]');
    const sequenceResumeButton = fixture.nativeElement.querySelector('button[uitestid=sequenceResumeButton]');
    const sequenceAbortButton = fixture.nativeElement.querySelector('button[uitestid=sequenceAbortButton]');

    // then
    expect(component.sequence.isFinished()).toBe(false);
    expect(component.sequence.isPaused()).toBe(true);

    expect(sequencePauseButton).toBeFalsy();
    expect(sequenceResumeButton).toBeTruthy();
    expect(sequenceAbortButton).toBeTruthy();

    expect(sequenceResumeButton.disabled).toBeFalsy();
    expect(sequenceAbortButton.disabled).toBeFalsy();
  });

  it('should resume sequence on click', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[3];
    fixture.detectChanges();

    const sequenceResumeButton = fixture.nativeElement.querySelector('button[uitestid=sequenceResumeButton]');
    const spyResume = jest.spyOn(component, 'triggerResumeSequence');

    // when
    sequenceResumeButton.dispatchEvent(new Event('click'));

    // then
    expect(spyResume).toHaveBeenCalled();
  });

  it('buttons should be disabled for finished sequence', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[0];
    fixture.detectChanges();

    const sequencePauseButton = fixture.nativeElement.querySelector('button[uitestid=sequencePauseButton]');
    const sequenceResumeButton = fixture.nativeElement.querySelector('button[uitestid=sequenceResumeButton]');
    const sequenceAbortButton = fixture.nativeElement.querySelector('button[uitestid=sequenceAbortButton]');

    // then
    expect(component.sequence.isFinished()).toBe(true);
    expect(component.sequence.isPaused()).toBe(false);

    expect(sequencePauseButton).toBeTruthy();
    expect(sequenceResumeButton).toBeFalsy();
    expect(sequenceAbortButton).toBeTruthy();

    expect(sequencePauseButton.disabled).toBeTruthy();
    expect(sequenceAbortButton.disabled).toBeTruthy();
  });
});
