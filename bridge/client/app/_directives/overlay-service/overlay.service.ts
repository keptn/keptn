import { ElementRef, Injectable } from '@angular/core';
import { Overlay, OverlayPositionBuilder, OverlayRef } from '@angular/cdk/overlay';
import { NavigationStart, Router } from '@angular/router';
import { filter, takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class OverlayService {
  constructor(
    private overlay: Overlay,
    private overlayPositionBuilder: OverlayPositionBuilder,
    private router: Router
  ) {}

  public initOverlay(
    width: string,
    height: string,
    hasBackdrop: boolean,
    elementRef: ElementRef,
    closeCallback: () => void
  ): OverlayRef {
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

    const overlayRef = this.overlay.create({
      positionStrategy,
      width,
      height,
      hasBackdrop,
      backdropClass: 'cdk-overlay-transparent-backdrop',
    });

    overlayRef.backdropClick().subscribe(() => {
      closeCallback();
    });

    return overlayRef;
  }

  public registerNavigationEvent(unsubscribe$: Subject<void>, closeCallback: () => void): void {
    // Close when navigation happens - to keep the overlay on the UI
    this.router.events
      .pipe(
        takeUntil(unsubscribe$),
        filter((event) => event instanceof NavigationStart)
      )
      .subscribe(() => {
        closeCallback();
      });
  }

  public closeOverlay(overlayRef: OverlayRef | undefined, elementRef: ElementRef): void {
    elementRef.nativeElement.disabled = false;
    overlayRef?.detach();
  }
}
