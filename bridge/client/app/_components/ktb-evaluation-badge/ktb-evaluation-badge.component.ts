import { Component, Input, NgZone, OnDestroy, TemplateRef, ViewChild } from '@angular/core';
import { Trace } from '../../_models/trace';
import { takeUntil } from 'rxjs/operators';
import { DtOverlay, DtOverlayConfig, DtOverlayRef } from '@dynatrace/barista-components/overlay';
import { Subject, Subscription } from 'rxjs';
import { EvaluationBadgeVariant, IEvaluationBadgeState } from './ktb-evaluation-badge.utils';

@Component({
  selector: 'ktb-evaluation-badge',
  templateUrl: './ktb-evaluation-badge.component.html',
  styleUrls: ['./ktb-evaluation-badge.component.scss'],
})
export class KtbEvaluationBadgeComponent implements OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private overlayRef?: DtOverlayRef<unknown>;
  private updateOverlayPositionSubscription = Subscription.EMPTY;
  public TraceClass = Trace;
  public EvaluationBadgeFillState = EvaluationBadgeVariant;
  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
  };

  @ViewChild('overlay', { static: true, read: TemplateRef }) overlayTemplate?: TemplateRef<unknown>;

  @Input() overlayDisabled = true;
  @Input() loading = false;
  @Input() evaluation?: Trace;
  @Input() evaluationState: Partial<IEvaluationBadgeState> = {
    isError: false,
    isSuccess: false,
    isWarning: false,
    fillState: EvaluationBadgeVariant.FILL,
  };

  constructor(private ngZone: NgZone, private _dtOverlay: DtOverlay) {}

  public showEvaluationOverlay(event: MouseEvent, data?: Trace): void {
    if (!this.overlayDisabled && this.overlayTemplate && data) {
      this.overlayRef = this._dtOverlay.create(event, this.overlayTemplate, { ...this.overlayConfig, data });
      this.updateEvaluationOverlayPosition();
    }
  }

  private updateEvaluationOverlayPosition(): void {
    this.updateOverlayPositionSubscription = this.ngZone.onMicrotaskEmpty
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this.overlayRef?.updatePosition();
        // if the content of the overlay changed after initialization the position stayed the same
      });
  }

  public hideEvaluationOverlay(): void {
    if (this.overlayRef) {
      this._dtOverlay.dismiss();
      this.updateOverlayPositionSubscription.unsubscribe();
      this.overlayRef = undefined;
    }
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
