import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbApprovalItemComponent } from './ktb-approval-item.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbApprovalItemModule } from './ktb-approval-item.module';
import { Trace } from '../../_models/trace';
import { DataService } from '../../_services/data.service';
import { EventTypes } from '../../../../shared/interfaces/event-types';
import { KeptnService } from '../../../../shared/models/keptn-service';
import { of } from 'rxjs';

describe('KtbEventItemComponent', () => {
  let component: KtbApprovalItemComponent;
  let fixture: ComponentFixture<KtbApprovalItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbApprovalItemModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbApprovalItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load an evaluation after the event is changed', () => {
    // given
    const dataService = TestBed.inject(DataService);
    const getEvaluationSpy = jest.spyOn(dataService, 'getTracesByContext');
    getEvaluationSpy.mockReturnValue(of([getEvaluationTrace()]));

    // when
    component.event = getTrace();

    // then
    expect(getEvaluationSpy).toHaveBeenCalledWith(
      'my-context',
      EventTypes.EVALUATION_FINISHED,
      KeptnService.LIGHTHOUSE_SERVICE,
      'dev',
      1
    );
    expect(component.evaluation).toEqual(getEvaluationTrace());
    expect(component.evaluationExists).toBe(true);
  });

  it('should not find an evaluation', () => {
    // given
    const dataService = TestBed.inject(DataService);
    jest.spyOn(dataService, 'getTracesByContext').mockReturnValue(of([]));

    // when
    component.event = getTrace();

    // then
    expect(component.evaluation).toBeUndefined();
    expect(component.evaluationExists).toBe(false);
  });

  it('should send and emit an accepted approval', () => {
    // given
    const dataService = TestBed.inject(DataService);
    const emitSpy = jest.spyOn(component.approvalSent, 'emit');
    const sendApprovalSpy = jest.spyOn(dataService, 'sendApprovalEvent');
    sendApprovalSpy.mockReturnValue(of({}));

    // when
    component.handleApproval(getTrace(), true);

    // then
    expect(sendApprovalSpy).toHaveBeenCalledWith(getTrace(), true);
    expect(emitSpy).toHaveBeenCalled();
    expect(component.approvalResult).toBe(true);
  });

  it('should send and emit a declined approval', () => {
    // given
    const dataService = TestBed.inject(DataService);
    const emitSpy = jest.spyOn(component.approvalSent, 'emit');
    const sendApprovalSpy = jest.spyOn(dataService, 'sendApprovalEvent');
    sendApprovalSpy.mockReturnValue(of({}));

    // when
    component.handleApproval(getTrace(), false);

    // then
    expect(sendApprovalSpy).toHaveBeenCalledWith(getTrace(), false);
    expect(emitSpy).toHaveBeenCalled();
    expect(component.approvalResult).toBe(false);
  });

  function getTrace(): Trace {
    return Trace.fromJSON({
      data: {
        project: 'sockshop',
        stage: 'dev',
        service: 'carts',
      },
      shkeptncontext: 'my-context',
      id: 'my-id',
      type: EventTypes.APPROVAL_STARTED,
    });
  }

  function getEvaluationTrace(): Trace {
    return Trace.fromJSON({
      data: {
        project: 'sockshop',
        stage: 'dev',
        service: 'carts',
      },
      shkeptncontext: 'my-context',
      id: 'my-id',
      type: EventTypes.EVALUATION_FINISHED,
    });
  }
});
