import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { HttpStateService } from '../../_services/http-state.service';
import { HttpProgressState, HttpState } from '../../_models/http-progress-state';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'ktb-http-loading-bar',
  templateUrl: './ktb-http-loading-bar.component.html',
  styleUrls: []
})
export class KtbHttpLoadingBarComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();

  private hideLoadingTimer?: ReturnType<typeof setTimeout>;
  private animateLoadingBarInterval?: ReturnType<typeof setInterval>;

  public loading = 0;
  @Input() public filterBy: string | null = null;

  public value = 0;
  public align: 'start' | 'end' = 'start';
  public state = 'recovered';

  private loadedUrls: string[] = [];
  private finishedUrls: string[] = [];

  constructor(private httpStateService: HttpStateService) {
  }

  ngOnInit(): void {
    this.httpStateService.state
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((progress: HttpState) => {
        if (progress && progress.url) {
          if (!this.filterBy || progress.url.indexOf(this.filterBy) !== -1) {
            if (progress.state === HttpProgressState.start) {
              this.showLoadingBar(progress.url);
            } else {
              this.hideLoadingBar(progress.url);
            }
          }
        }
      });
  }

  showLoadingBar(url: string): void {
    if (!this.loadedUrls.includes(url)) {
      this.loadedUrls.push(url);
      if (this.loading === 0) {
        this.animateLoadingBarInterval = setInterval(() => this.animateLoadingBar(), 500);
      }
      this.loading++;
    }
  }

  hideLoadingBar(url: string): void {
    if (!this.finishedUrls.includes(url)) {
      this.finishedUrls.push(url);
      this.hideLoadingTimer = setTimeout(() => this.resetValues(), 500);
    }
  }

  isLoading(): boolean {
    return this.loading > 0;
  }

  resetValues(): void {
    if (this.loading > 0) {
      this.loading--;
    }
    if (this.loading === 0) {
      if (this.animateLoadingBarInterval) {
        clearInterval(this.animateLoadingBarInterval);
      }
      this.value = 0;
      this.align = 'start';
    }
  }

  animateLoadingBar(): void {
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
