import { ElementRef, Injectable } from '@angular/core';
import { Overlay, OverlayPositionBuilder, OverlayRef } from '@angular/cdk/overlay';

@Injectable({
  providedIn: 'root',
})
export class OverlayService {
  constructor(private overlay: Overlay, private overlayPositionBuilder: OverlayPositionBuilder) {}

  public initOverlay(width: string, height: string, hasBackdrop: boolean, elementRef: ElementRef): OverlayRef {
    const positionStrategy = this.overlayPositionBuilder.flexibleConnectedTo(elementRef).withPositions([
      {
        originX: 'start',
        originY: 'bottom',
        overlayX: 'start',
        overlayY: 'top',
        offsetY: 10,
        offsetX: -20,
      },
    ]);

    return this.overlay.create({
      positionStrategy,
      width,
      height,
      hasBackdrop,
      backdropClass: 'cdk-overlay-transparent-backdrop',
    });
  }

  public closeOverlay(overlayRef: OverlayRef | undefined, elementRef: ElementRef): void {
    elementRef.nativeElement.disabled = false;
    overlayRef?.detach();
  }
}
