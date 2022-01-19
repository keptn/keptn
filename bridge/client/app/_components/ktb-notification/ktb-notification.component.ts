import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ComponentFactoryResolver,
  EventEmitter,
  HostListener,
  Input,
  OnDestroy,
  Output,
  ViewChild,
  ViewContainerRef,
} from '@angular/core';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { Location } from '@angular/common';
import { ComponentInfo, NotificationType } from '../../_models/notification';

const defaultTimeMs = 5_000;

@Component({
  selector: 'ktb-notification[message]',
  templateUrl: './ktb-notification.component.html',
  styleUrls: ['./ktb-notification.component.scss'],
  animations: [
    trigger('fade', [
      state('in', style({ opacity: 1 })),
      state('out', style({ opacity: 0 })),
      transition('in => out', animate('{{ duration }}ms')),
    ]),
  ],
})
export class KtbNotificationComponent implements AfterViewInit, OnDestroy {
  private _timeoutMs = defaultTimeMs;
  private timeout?: ReturnType<typeof setTimeout>;
  private fadeOutDelay?: ReturnType<typeof setTimeout>;
  public fadeStatus: 'in' | 'out' = 'in';
  public fadeOutDuration = 3_000;
  public NotificationType = NotificationType;

  @Input() severity: NotificationType = NotificationType.SUCCESS;
  @Input() message = '';

  /**
   * timeout for the notification in milliseconds.
   * <br>If -1, the timeout is disabled.
   * <br>If not provided, the value is set to a default one
   * @param time
   */
  @Input()
  set timeoutMs(time: number | undefined) {
    this._timeoutMs = time ?? defaultTimeMs;
  }
  get timeoutMs(): number | undefined {
    return this._timeoutMs;
  }
  @Input() componentInfo?: ComponentInfo;
  @Output() close: EventEmitter<void> = new EventEmitter<void>();
  @ViewChild('showComponent', { read: ViewContainerRef }) componentViewRef?: ViewContainerRef;

  constructor(
    private _changeDetectorRef: ChangeDetectorRef,
    private resolver: ComponentFactoryResolver,
    public location: Location /*used for create project link*/
  ) {}

  @HostListener('mouseover')
  public onHover(): void {
    this.stopTimeout();
    this.fadeStatus = 'in';
  }

  @HostListener('mouseleave')
  public onLeave(): void {
    this.startTimeout();
  }

  public ngAfterViewInit(): void {
    this.loadComponent();
    this.startTimeout();
    this._changeDetectorRef.detectChanges();
  }

  private loadComponent(): void {
    if (this.componentInfo && this.componentViewRef) {
      const viewContainerRef = this.componentViewRef;
      viewContainerRef.clear();

      const factory = this.resolver.resolveComponentFactory(this.componentInfo.component);
      const componentRef = viewContainerRef.createComponent(factory);
      (componentRef.instance as { data: Record<string, unknown> }).data = this.componentInfo.data;
    }
  }

  public closeComponent(): void {
    this.close.next();
  }

  private startTimeout(): void {
    if (this.timeoutMs !== -1) {
      //Only start fade out the last 3 seconds
      this.fadeOutDelay = setTimeout(() => {
        this.fadeStatus = 'out';
      }, this._timeoutMs - this.fadeOutDuration);
      this.timeout = setTimeout(() => {
        this.closeComponent();
      }, this._timeoutMs);
    }
  }

  private stopTimeout(): void {
    clearTimeout(this.timeout);
    clearTimeout(this.fadeOutDelay);
  }

  public ngOnDestroy(): void {
    this.stopTimeout();
  }
}
