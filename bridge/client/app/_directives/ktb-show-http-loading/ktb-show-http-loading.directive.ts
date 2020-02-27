import {Directive, Input, OnInit, TemplateRef, ViewContainerRef} from '@angular/core';
import {HttpStateService} from "../../_services/http-state.service";
import {HttpProgressState, HttpState} from "../../_models/http-progress-state";

@Directive({
  selector: '[ktbShowHttpLoading]'
})
export class KtbShowHttpLoadingDirective implements OnInit {

  public filterBy: string | null = null;
  private hideTimer;

  @Input() set ktbShowHttpLoading(filterBy: string) {
    this.filterBy = filterBy;
  }

  constructor(private httpStateService: HttpStateService, private templateRef: TemplateRef<any>, private viewContainer: ViewContainerRef) { }

  ngOnInit(): void {
    this.httpStateService.state.subscribe((progress: HttpState) => {
      if (progress && progress.url) {
        if(!this.filterBy || progress.url.indexOf(this.filterBy) !== -1) {
          if(progress.state === HttpProgressState.start) {
            this.showElement();
          } else {
            this.hideElement();
          }
        }
      }
    });
  }

  showElement() {
    clearTimeout(this.hideTimer);
    this.viewContainer.createEmbeddedView(this.templateRef);
  }

  hideElement() {
    this.hideTimer = setTimeout(() => {
      this.viewContainer.clear();
    }, 500);
  }

}
