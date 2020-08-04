import {Directive, Input, OnDestroy, OnInit, TemplateRef, ViewContainerRef} from '@angular/core';
import {HttpStateService} from "../../_services/http-state.service";
import {HttpProgressState, HttpState} from "../../_models/http-progress-state";
import {Subject} from "rxjs";
import {takeUntil} from "rxjs/operators";

@Directive({
  selector: '[ktbHideHttpLoading]'
})
export class KtbHideHttpLoadingDirective implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();

  public filterBy: string | null = null;
  private showTimer;

  @Input() set ktbHideHttpLoading(filterBy: string) {
    this.filterBy = filterBy;
  }

  constructor(private httpStateService: HttpStateService, private templateRef: TemplateRef<any>, private viewContainer: ViewContainerRef) { }

  ngOnInit(): void {
    this.httpStateService.state
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((progress: HttpState) => {
        if (progress && progress.url) {
          if(!this.filterBy || progress.url.indexOf(this.filterBy) !== -1) {
            if(progress.state === HttpProgressState.start) {
              this.hideElement();
            } else {
              this.showElement();
            }
          }
        }
      });
  }

  showElement() {
    this.showTimer = setTimeout(() => {
      this.viewContainer.clear();
      this.viewContainer.createEmbeddedView(this.templateRef);
    }, 500);
  }

  hideElement() {
    clearTimeout(this.showTimer);
    this.viewContainer.clear();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
