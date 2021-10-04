import { Directive, Input, OnDestroy, OnInit, TemplateRef, ViewContainerRef } from '@angular/core';
import { HttpStateService } from '../../_services/http-state.service';
import { HttpProgressState, HttpState } from '../../_models/http-progress-state';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Directive({
  selector: '[ktbShowHttpLoading]',
})
export class KtbShowHttpLoadingDirective implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  public filterBy: string | null = null;
  private hideTimer?: ReturnType<typeof setTimeout>;

  @Input() set ktbShowHttpLoading(filterBy: string) {
    this.filterBy = filterBy;
  }

  constructor(
    private httpStateService: HttpStateService,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef
  ) {}

  ngOnInit(): void {
    this.httpStateService.state.pipe(takeUntil(this.unsubscribe$)).subscribe((progress: HttpState) => {
      if (progress && progress.url) {
        if (!this.filterBy || progress.url.indexOf(this.filterBy) !== -1) {
          if (progress.state === HttpProgressState.START) {
            this.showElement();
          } else {
            this.hideElement();
          }
        }
      }
    });
  }

  showElement(): void {
    if (this.hideTimer) {
      clearTimeout(this.hideTimer);
    }
    this.viewContainer.createEmbeddedView(this.templateRef);
  }

  hideElement(): void {
    this.hideTimer = setTimeout(() => {
      this.viewContainer.clear();
    }, 500);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
