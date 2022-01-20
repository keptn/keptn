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
import { Notification, NotificationType } from '../../_models/notification';

@Component({
  selector: 'ktb-notification[notification]',
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
  private hideTimeout?: ReturnType<typeof setTimeout>;
  private fadeOutDelay?: ReturnType<typeof setTimeout>;
  public fadeStatus: 'in' | 'out' = 'in';
  public fadeOutDuration = 3_000;
  public NotificationType = NotificationType;

  @Input() notification!: Notification;
  @Output() hide: EventEmitter<void> = new EventEmitter<void>();
  @ViewChild('showComponent', { read: ViewContainerRef }) componentViewRef?: ViewContainerRef;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private resolver: ComponentFactoryResolver) {}

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
    if (this.notification.componentInfo && this.componentViewRef) {
      this.componentViewRef.clear();

      const factory = this.resolver.resolveComponentFactory(this.notification.componentInfo.component);
      const componentRef = this.componentViewRef.createComponent(factory);
      Object.assign(componentRef.instance, this.notification.componentInfo.data);
    }
  }

  public closeComponent(): void {
    this.hide.next();
  }

  private startTimeout(): void {
    if (!this.notification.isPinned) {
      //Only start fade out the last 3 seconds
      this.fadeOutDelay = setTimeout(() => {
        this.fadeStatus = 'out';
      }, this.notification.time - this.fadeOutDuration);

      this.hideTimeout = setTimeout(() => {
        this.closeComponent();
      }, this.notification.time);
    }
  }

  private stopTimeout(): void {
    clearTimeout(this.hideTimeout);
    clearTimeout(this.fadeOutDelay);
  }

  public ngOnDestroy(): void {
    this.stopTimeout();
  }
}
