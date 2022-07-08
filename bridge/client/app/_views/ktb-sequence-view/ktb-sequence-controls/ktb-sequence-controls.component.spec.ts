import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceControlsComponent } from './ktb-sequence-controls.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';
import { ApiService } from '../../../_services/api.service';
import { ApiServiceMock } from '../../../_services/api.service.mock';
import { SequencesMock } from '../../../_services/_mockData/sequences.mock';
import { KtbSequenceViewModule } from '../ktb-sequence-view.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('KtbSequenceControlsComponent', () => {
  let component: KtbSequenceControlsComponent;
  let fixture: ComponentFixture<KtbSequenceControlsComponent>;
  const projectName = 'sockshop';

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSequenceViewModule, BrowserAnimationsModule, HttpClientTestingModule],
      providers: [
        {
          provide: ApiService,
          useClass: ApiServiceMock,
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

    fixture = TestBed.createComponent(KtbSequenceControlsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show pause and abort for started sequence', () => {
    // given
    component.sequence = SequencesMock[7];
    fixture.detectChanges();

    const sequencePauseButton = fixture.nativeElement.querySelector('button[uitestid=sequencePauseButton]');
    const sequenceResumeButton = fixture.nativeElement.querySelector('button[uitestid=sequenceResumeButton]');
    const sequenceAbortButton = fixture.nativeElement.querySelector('button[uitestid=sequenceAbortButton]');

    // then
    expect(component.sequence?.isFinished()).toBe(false);
    expect(component.sequence?.isPaused()).toBe(false);

    expect(sequencePauseButton).toBeTruthy();
    expect(sequenceResumeButton).toBeFalsy();
    expect(sequenceAbortButton).toBeTruthy();

    expect(sequencePauseButton.disabled).toBeFalsy();
    expect(sequenceAbortButton.disabled).toBeFalsy();
  });

  it('should pause sequence on click', () => {
    // given
    component.sequence = SequencesMock[7];
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
    component.sequence = SequencesMock[7];
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
    component.sequence = SequencesMock[4];
    fixture.detectChanges();

    const sequencePauseButton = fixture.nativeElement.querySelector('button[uitestid=sequencePauseButton]');
    const sequenceResumeButton = fixture.nativeElement.querySelector('button[uitestid=sequenceResumeButton]');
    const sequenceAbortButton = fixture.nativeElement.querySelector('button[uitestid=sequenceAbortButton]');

    // then
    expect(component.sequence?.isFinished()).toBe(false);
    expect(component.sequence?.isPaused()).toBe(true);

    expect(sequencePauseButton).toBeFalsy();
    expect(sequenceResumeButton).toBeTruthy();
    expect(sequenceAbortButton).toBeTruthy();

    expect(sequenceResumeButton.disabled).toBeFalsy();
    expect(sequenceAbortButton.disabled).toBeFalsy();
  });

  it('should resume sequence on click', () => {
    // given
    component.sequence = SequencesMock[4];
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
    component.sequence = SequencesMock[1];
    fixture.detectChanges();

    const sequencePauseButton = fixture.nativeElement.querySelector('button[uitestid=sequencePauseButton]');
    const sequenceResumeButton = fixture.nativeElement.querySelector('button[uitestid=sequenceResumeButton]');
    const sequenceAbortButton = fixture.nativeElement.querySelector('button[uitestid=sequenceAbortButton]');

    // then
    expect(component.sequence?.isFinished()).toBe(true);
    expect(component.sequence?.isPaused()).toBe(false);

    expect(sequencePauseButton).toBeTruthy();
    expect(sequenceResumeButton).toBeFalsy();
    expect(sequenceAbortButton).toBeTruthy();

    expect(sequencePauseButton.disabled).toBeTruthy();
    expect(sequenceAbortButton.disabled).toBeTruthy();
  });
});
