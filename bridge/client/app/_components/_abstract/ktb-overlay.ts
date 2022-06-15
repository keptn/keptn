import { ElementRef } from '@angular/core';
import { OverlayService } from '../../_directives/overlay-service/overlay.service';
import { Subject } from 'rxjs';
import { OverlayRef } from '@angular/cdk/overlay';

export abstract class KtbOverlay {
  protected elementRef: ElementRef;
  protected overlayRef?: OverlayRef;
  protected overlayService: OverlayService;
  protected unsubscribe$: Subject<void> = new Subject();
  private readonly closeCallback = (): void => this.close();
  private readonly width: string;
  private readonly height: string;

  protected constructor(elementRef: ElementRef, overlayService: OverlayService, width: string, height: string) {
    this.elementRef = elementRef;
    this.overlayService = overlayService;
    this.width = width;
    this.height = height;
    // Close when navigation happens - to keep the overlay on the UI
    this.overlayService.registerNavigationEvent(this.unsubscribe$, this.closeCallback);
  }

  public onInit(): void {
    this.overlayRef = this.overlayService.initOverlay(
      this.width,
      this.height,
      true,
      this.elementRef,
      this.closeCallback
    );
  }

  public onDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public close(): void {
    this.overlayService.closeOverlay(this.overlayRef, this.elementRef);
  }
}
