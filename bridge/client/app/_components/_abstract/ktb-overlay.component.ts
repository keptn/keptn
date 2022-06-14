import { ElementRef, OnDestroy, OnInit } from '@angular/core';
import { OverlayService } from '../../_directives/overlay-service/overlay.service';
import { Subject } from 'rxjs';
import { OverlayRef } from '@angular/cdk/overlay';

export abstract class KtbOverlayComponent implements OnInit, OnDestroy {
  protected overlayRef?: OverlayRef;
  protected unsubscribe$: Subject<void> = new Subject();

  protected constructor(
    protected elementRef: ElementRef,
    protected overlayService: OverlayService,
    protected width: string,
    protected height: string
  ) {
    // Close when navigation happens - to keep the overlay on the UI
    this.overlayService.registerNavigationEvent(this.unsubscribe$, this.close.bind(this));
  }

  public ngOnInit(): void {
    const closeCallback = (): void => this.close();
    this.overlayRef = this.overlayService.initOverlay(this.width, this.height, true, this.elementRef, closeCallback);
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public close(): void {
    this.overlayService.closeOverlay(this.overlayRef, this.elementRef);
  }
}
