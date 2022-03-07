import {
  Component,
  ComponentRef,
  Directive,
  ElementRef,
  EventEmitter,
  HostListener,
  Input,
  OnDestroy,
  OnInit,
  Output,
} from '@angular/core';
import { OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';
import { Router } from '@angular/router';
import moment from 'moment';
import { Timeframe } from '../../_models/timeframe';
import { OverlayService } from '../../_directives/overlay-service/overlay.service';
import { Subject } from 'rxjs';

@Directive({
  selector: '[ktbDatetimePicker]',
})
export class KtbDatetimePickerDirective implements OnInit, OnDestroy {
  private overlayRef?: OverlayRef;
  private contentRef: ComponentRef<KtbDatetimePickerComponent> | undefined;
  private unsubscribe$: Subject<void> = new Subject();

  @Input() timeEnabled = false;
  @Input() secondsEnabled = false;
  @Output() selectedDateTime: EventEmitter<string> = new EventEmitter<string>();

  @HostListener('click')
  show(): void {
    const dateTimePickerPortal: ComponentPortal<KtbDatetimePickerComponent> = new ComponentPortal(
      // eslint-disable-next-line @typescript-eslint/no-use-before-define
      KtbDatetimePickerComponent
    );
    // Disable origin to prevent 'Host has already a portal attached' error
    this.elementRef.nativeElement.disabled = true;

    this.contentRef = this.overlayRef?.attach(dateTimePickerPortal);
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

  constructor(private elementRef: ElementRef, private router: Router, private overlayService: OverlayService) {
    // Close when navigation happens - to keep the overlay on the UI
    this.overlayService.registerNavigationEvent(this.unsubscribe$, this.close.bind(this));
  }

  public ngOnInit(): void {
    this.overlayRef = this.overlayService.initOverlay('350px', '400px', true, this.elementRef, this.close.bind(this));
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public close(): void {
    this.overlayService.closeOverlay(this.overlayRef, this.elementRef);
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
  public maxTimeValues: Timeframe = {
    hours: 23,
    minutes: 59,
    seconds: 59,
    millis: 999,
    micros: 999,
  };
  private selectedDate = moment().hours(0).minutes(0).seconds(0).milliseconds(0);
  private selectedTime: Timeframe | undefined;

  public changeDate(event: Date): void {
    this.selectedDate = moment(event).hours(0).minutes(0).seconds(0).milliseconds(0);
  }

  public changeTime(time: Timeframe): void {
    this.selectedTime = time;

    if (this.secondsEnabled) {
      this.disabled = !(
        (time.hours !== undefined && time.minutes !== undefined && time.seconds !== undefined) ||
        (time.hours === undefined && time.minutes === undefined && time.seconds === undefined)
      );
    } else {
      this.disabled = !(
        (time.hours !== undefined && time.minutes !== undefined) ||
        (time.hours === undefined && time.minutes === undefined)
      );
    }
  }

  public setDateTime(): void {
    if (
      this.selectedTime !== undefined &&
      this.selectedTime.hours !== undefined &&
      this.selectedTime.minutes !== undefined
    ) {
      this.selectedDate.hours(this.selectedTime.hours);
      this.selectedDate.minutes(this.selectedTime.minutes);
      if (this.secondsEnabled && this.selectedTime.seconds !== undefined) {
        this.selectedDate.seconds(this.selectedTime.seconds);
      }
    }
    this.selectedDateTime.emit(this.selectedDate.toISOString());
    this.closeDialog.emit();
  }
}
