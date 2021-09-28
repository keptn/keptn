import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { HttpStateService } from '../../_services/http-state.service';
import { HttpProgressState, HttpState } from '../../_models/http-progress-state';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'ktb-http-loading-bar',
  templateUrl: './ktb-http-loading-bar.component.html',
  styleUrls: [],
})
export class KtbHttpLoadingBarComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  private hideLoadingTimer?: ReturnType<typeof setTimeout>;
  private animateLoadingBarInterval?: ReturnType<typeof setInterval>;

  public loading = false;
  @Input() public filterBy: string | null = null;

  public value = 0;
  public align: 'start' | 'end' = 'start';
  public state = 'recovered';

  constructor(private httpStateService: HttpStateService) {}

  ngOnInit() {
    this.httpStateService.state.pipe(takeUntil(this.unsubscribe$)).subscribe((progress: HttpState) => {
      if (progress && progress.url) {
        if (!this.filterBy || progress.url.indexOf(this.filterBy) !== -1) {
          if (progress.state === HttpProgressState.start) {
            this.showLoadingBar();
          } else {
            this.hideLoadingBar();
          }
        }
      }
    });
  }

  showLoadingBar() {
    if (!this.loading) {
      if (this.hideLoadingTimer) {
        clearTimeout(this.hideLoadingTimer);
      }
      this.loading = true;
      this.animateLoadingBarInterval = setInterval(() => this.animateLoadingBar(), 500);
    }
  }

  hideLoadingBar() {
    if (this.loading) {
      if (this.animateLoadingBarInterval) {
        clearInterval(this.animateLoadingBarInterval);
      }
      this.hideLoadingTimer = setTimeout(() => this.resetValues(), 500);
    }
  }

  resetValues() {
    this.loading = false;
    this.value = 0;
    this.align = 'start';
  }

  animateLoadingBar() {
    if (this.align === 'start') {
      if (this.value < 100) {
        this.value = 100;
      } else {
        this.align = 'end';
      }
    } else {
      if (this.value > 0) {
        this.value = 0;
      } else {
        this.align = 'start';
      }
    }
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
