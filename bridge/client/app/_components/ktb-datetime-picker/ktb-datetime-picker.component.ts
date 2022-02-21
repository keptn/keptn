import {
  Component,
  ComponentRef,
  Directive,
  ElementRef,
  EventEmitter,
  HostListener,
  Input,
  OnInit,
  Output,
} from '@angular/core';
import { Overlay, OverlayPositionBuilder, OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';
import { NavigationStart, Router } from '@angular/router';
import { filter } from 'rxjs/operators';
import moment from 'moment';
import { Timeframe } from '../../_models/timeframe';

@Directive({
  selector: '[ktbDatetimePicker]',
})
export class KtbDatetimePickerDirective implements OnInit {
  private overlayRef?: OverlayRef;
  private contentRef: ComponentRef<KtbDatetimePickerComponent> | undefined;

  @Input() timeEnabled = false;
  @Input() secondsEnabled = false;
  @Output() selectedDateTime: EventEmitter<string> = new EventEmitter<string>();

  @HostListener('click')
  show(): void {
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    const tooltipPortal: ComponentPortal<KtbDatetimePickerComponent> = new ComponentPortal(KtbDatetimePickerComponent);
    // Disable origin to prevent 'Host has already a portal attached' error
    this.elementRef.nativeElement.disabled = true;

    this.contentRef = this.overlayRef?.attach(tooltipPortal);
    if (this.contentRef) {
      this.contentRef.instance.timeEnabled = this.timeEnabled;
      this.contentRef.instance.secondsEnabled = this.secondsEnabled;
      this.contentRef.instance.closeDialog.subscribe(() => {
        this.close();
      });

      this.contentRef.instance.selectedDateTime.subscribe((selected) => {
        this.selectedDateTime.emit(selected);
        this.close();
      });
    }
  }

  constructor(
    private overlay: Overlay,
    private overlayPositionBuilder: OverlayPositionBuilder,
    private elementRef: ElementRef,
    private router: Router
  ) {
    // Close when navigation happens - to keep the overlay on the UI
    this.router.events.pipe(filter((event) => event instanceof NavigationStart)).subscribe(() => {
      this.close();
    });
  }

  public ngOnInit(): void {
    const positionStrategy = this.overlayPositionBuilder.flexibleConnectedTo(this.elementRef).withPositions([
      {
        originX: 'start',
        originY: 'bottom',
        overlayX: 'start',
        overlayY: 'top',
        offsetY: 10,
        offsetX: -20,
      },
    ]);

    this.overlayRef = this.overlay.create({
      positionStrategy,
      width: '350px',
      height: '400px',
      hasBackdrop: true,
      backdropClass: 'cdk-overlay-transparent-backdrop',
    });
    this.overlayRef.backdropClick().subscribe(() => {
      this.close();
    });
  }

  public close(): void {
    this.elementRef.nativeElement.disabled = false;
    this.overlayRef?.detach();
  }
}

@Component({
  selector: 'ktb-datetime-picker',
  templateUrl: './ktb-datetime-picker.component.html',
  styleUrls: ['./ktb-datetime-picker.component.scss'],
})
export class KtbDatetimePickerComponent {
  @Input() timeEnabled = false;
  @Input() secondsEnabled = false;
  @Output() closeDialog: EventEmitter<void> = new EventEmitter<void>();
  @Output() selectedDateTime: EventEmitter<string> = new EventEmitter<string>();

  public disabled = false;
  public maxDate: Date = new Date();
  private selectedDate = moment();
  private selectedTime: Timeframe | undefined;

  public changeDate(event: Date): void {
    this.selectedDate = moment(event);
  }

  public changeTime(time: Timeframe): void {
    this.selectedTime = time;
    this.disabled =
      time.hours === undefined || time.minutes === undefined || (this.secondsEnabled && time.seconds === undefined);
  }

  public setDateTime(): void {
    if (
      this.selectedTime !== undefined &&
      this.selectedTime.hours !== undefined &&
      this.selectedTime.minutes !== undefined
    ) {
      this.selectedDate.set('hours', this.selectedTime.hours);
      this.selectedDate.set('minutes', this.selectedTime.minutes);
      if (this.secondsEnabled && this.selectedTime.seconds !== undefined) {
        this.selectedDate.set('seconds', this.selectedTime.seconds);
      }
    }

    this.selectedDateTime.emit(this.selectedDate.toISOString());
    this.closeDialog.emit();
  }
}
