import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEvaluationBadgeComponent } from './ktb-evaluation-badge.component';
import { Trace } from '../../_models/trace';
import { DtOverlay, DtOverlayRef } from '@dynatrace/barista-components/overlay';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { NgZone } from '@angular/core';
import { KtbEvaluationBadgeModule } from './ktb-evaluation-badge.module';

describe('KtbEvaluationBadgeComponent', () => {
  let component: KtbEvaluationBadgeComponent;
  let fixture: ComponentFixture<KtbEvaluationBadgeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BrowserAnimationsModule, KtbEvaluationBadgeModule, HttpClientTestingModule],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEvaluationBadgeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('showEvaluationOverlay', () => {
    it('should show overlay', () => {
      // given
      const overlay = TestBed.inject(DtOverlay);
      const createOverlaySpy = jest.spyOn(overlay, 'create');
      component.overlayDisabled = false;

      // when
      component.showEvaluationOverlay(new MouseEvent('enter'), Trace.fromJSON({}));

      // then
      expect(createOverlaySpy).toHaveBeenCalled();
    });

    it('should not show overlay if overlay is disabled', () => {
      // given
      const overlay = TestBed.inject(DtOverlay);
      const createOverlaySpy = jest.spyOn(overlay, 'create');
      component.overlayDisabled = true;

      // when
      component.showEvaluationOverlay(new MouseEvent('enter'), Trace.fromJSON({}));

      // then
      expect(createOverlaySpy).not.toHaveBeenCalled();
    });

    it('should not show overlay if evaluation is undefined', () => {
      // given
      const overlay = TestBed.inject(DtOverlay);
      const createOverlaySpy = jest.spyOn(overlay, 'create');
      component.overlayDisabled = false;

      // when
      component.showEvaluationOverlay(new MouseEvent('enter'), undefined);

      // then
      expect(createOverlaySpy).not.toHaveBeenCalled();
    });

    it('should not show overlay if template is not found', () => {
      // given
      const overlay = TestBed.inject(DtOverlay);
      const createOverlaySpy = jest.spyOn(overlay, 'create');
      component.overlayDisabled = true;
      component.overlayTemplate = undefined;

      // when
      component.showEvaluationOverlay(new MouseEvent('enter'), Trace.fromJSON({}));

      // then
      expect(createOverlaySpy).not.toHaveBeenCalled();
    });
  });

  describe('hideEvaluationOverlay', () => {
    it('should dismiss overlay if it exists', () => {
      // given
      const overlay = TestBed.inject(DtOverlay);
      const dismissOverlaySpy = jest.spyOn(overlay, 'dismiss');
      component.overlayDisabled = false;
      component.showEvaluationOverlay(new MouseEvent('enter'), Trace.fromJSON({}));

      // when
      component.hideEvaluationOverlay();

      // then
      expect(dismissOverlaySpy).toHaveBeenCalled();
    });

    it('should not dismiss overlay if it does not exists', () => {
      // given
      const overlay = TestBed.inject(DtOverlay);
      const dismissOverlaySpy = jest.spyOn(overlay, 'dismiss');

      // when
      component.hideEvaluationOverlay();

      // then
      expect(dismissOverlaySpy).not.toHaveBeenCalled();
    });
  });

  describe('updateEvaluationOverlayPosition', () => {
    it('should update overlay position after creation', () => {
      // given
      const overlay = TestBed.inject(DtOverlay);
      const updatePositionSpy = jest.fn();
      jest
        .spyOn(overlay, 'create')
        .mockReturnValue({ updatePosition: updatePositionSpy } as unknown as DtOverlayRef<unknown>);
      component.overlayDisabled = false;
      component.showEvaluationOverlay(new MouseEvent('enter'), Trace.fromJSON({}));

      // when
      TestBed.inject(NgZone).onMicrotaskEmpty.emit('next');

      // then
      expect(updatePositionSpy).toHaveBeenCalled();
    });
  });
});
