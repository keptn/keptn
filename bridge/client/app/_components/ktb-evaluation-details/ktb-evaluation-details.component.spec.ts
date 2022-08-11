import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEvaluationDetailsComponent } from './ktb-evaluation-details.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { KtbEvaluationDetailsModule } from './ktb-evaluation-details.module';
import { ElementRef, EmbeddedViewRef, TemplateRef } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ClipboardService } from '../../_services/clipboard.service';
import { Trace } from '../../_models/trace';
import { EvaluationsMock } from '../../_services/_mockData/evaluations.mock';
import { DataService } from '../../_services/data.service';

describe(KtbEvaluationDetailsComponent.name, () => {
  let component: KtbEvaluationDetailsComponent;
  let fixture: ComponentFixture<KtbEvaluationDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbEvaluationDetailsModule, BrowserAnimationsModule, HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEvaluationDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('SLO dialog', () => {
    let openDialogSpy: jest.SpyInstance;
    beforeEach(() => {
      openDialogSpy = jest
        .spyOn(TestBed.inject(MatDialog), 'open')
        .mockReturnValue({ close: () => {} } as MatDialogRef<string>);
    });

    it('should show SLO dialog', () => {
      // given, when
      const template = showSloDialog('bXlTTE8=');

      // then
      expect(openDialogSpy).toHaveBeenCalledWith(template, {
        data: 'mySLO',
      });
      expect(Reflect.get(component, 'sloDialogRef')).not.toBeUndefined();
    });

    it('should close SLO dialog', () => {
      // given
      showSloDialog('bXlTTE8=');
      const closeSpy = jest.spyOn(Reflect.get(component, 'sloDialogRef'), 'close');

      // when
      component.closeSloDialog();

      // then
      expect(closeSpy).toHaveBeenCalled();
    });
  });

  it('should copy payload to clipboard', async () => {
    // given
    const clipboardSpy = jest.spyOn(TestBed.inject(ClipboardService), 'copy');
    document.execCommand = jest.fn();

    // when
    component.copySloPayload('myPayload');

    // then
    expect(document.execCommand).toHaveBeenCalledWith('copy');
    expect(clipboardSpy).toHaveBeenCalledWith('myPayload', 'slo payload');
  });

  describe('invalidate dialog', () => {
    let openDialogSpy: jest.SpyInstance;
    beforeEach(() => {
      openDialogSpy = jest
        .spyOn(TestBed.inject(MatDialog), 'open')
        .mockReturnValue({ close: () => {} } as MatDialogRef<Trace | undefined>);
    });

    it('should show invalidate dialog', () => {
      // given
      const trace = EvaluationsMock;
      component.evaluationData = {
        evaluation: trace,
        shouldSelect: true,
      };

      // when
      const template = showInvalidateDialog();

      // then
      expect(openDialogSpy).toHaveBeenCalledWith(template, {
        data: trace,
      });
      expect(Reflect.get(component, 'invalidateEvaluationDialogRef')).not.toBeUndefined();
    });

    it('should close invalidate dialog', () => {
      // given
      showInvalidateDialog();
      const closeSpy = jest.spyOn(Reflect.get(component, 'invalidateEvaluationDialogRef'), 'close');

      // when
      component.closeInvalidateEvaluationDialog();

      // then
      expect(closeSpy).toHaveBeenCalled();
    });

    it('should close invalidate dialog after invalidating', () => {
      // given
      const trace = EvaluationsMock;
      component.evaluationData = {
        evaluation: trace,
        shouldSelect: true,
      };
      showInvalidateDialog();
      const closeSpy = jest.spyOn(Reflect.get(component, 'invalidateEvaluationDialogRef'), 'close');
      const invalidateSpy = jest.spyOn(TestBed.inject(DataService), 'invalidateEvaluation');

      // when
      component.invalidateEvaluation(trace, 'myReason');

      // then
      expect(closeSpy).toHaveBeenCalled();
      expect(invalidateSpy).toHaveBeenCalledWith(trace, 'myReason');
    });
  });

  it('should map evaluation to evaluationData', () => {
    // given, when
    component.evaluation = EvaluationsMock;

    // then
    expect(component.evaluationData).toEqual({
      evaluation: EvaluationsMock,
      shouldSelect: true,
    });
    expect(component.evaluation).toEqual(EvaluationsMock);
  });

  function showSloDialog(content: string): TemplateRef<string> {
    const template = getTemplateRefMock();
    component.showSloDialog(content, template);
    return template;
  }

  function showInvalidateDialog(): TemplateRef<Trace | undefined> {
    const template = {
      get elementRef(): ElementRef {
        return {
          nativeElement: undefined,
        };
      },
      createEmbeddedView(): EmbeddedViewRef<Trace | undefined> {
        return undefined as unknown as EmbeddedViewRef<Trace | undefined>;
      },
    };
    component.invalidateEvaluationTrigger(template);
    return template;
  }

  function getTemplateRefMock(): TemplateRef<string> {
    return {
      get elementRef(): ElementRef {
        return {
          nativeElement: undefined,
        };
      },
      createEmbeddedView(): EmbeddedViewRef<string> {
        return undefined as unknown as EmbeddedViewRef<string>;
      },
    };
  }
});
