import {Component, Input, OnInit} from '@angular/core';
import {HttpStateService} from "../../_services/http-state.service";
import {HttpProgressState, HttpState} from "../../_models/http-progress-state";

@Component({
  selector: 'ktb-http-loading-bar',
  templateUrl: './ktb-http-loading-bar.component.html',
  styleUrls: ['./ktb-http-loading-bar.component.scss']
})
export class KtbHttpLoadingBarComponent implements OnInit {

  public loading = false;
  @Input() public filterBy: string | null = null;

  private hideLoadingTimer;
  private animateLoadingBarInterval;

  private value = 0;
  private align = 'start';
  private state = 'recovered';

  constructor(private httpStateService: HttpStateService) { }

  ngOnInit() {
    this.httpStateService.state.subscribe((progress: HttpState) => {
      if (progress && progress.url) {
        if(!this.filterBy || progress.url.indexOf(this.filterBy) !== -1) {
          if(progress.state === HttpProgressState.start)
            this.showLoadingBar();
          else
            this.hideLoadingBar();
        }
      }
    });
  }

  showLoadingBar() {
    if(this.loading && !this.hideLoadingTimer)
      return;
    clearTimeout(this.hideLoadingTimer);
    this.loading = true;
    this.animateLoadingBarInterval = setInterval(() => this.animateLoadingBar(), 500);
  }

  hideLoadingBar() {
    if(!this.loading)
      return;
    clearInterval(this.animateLoadingBarInterval);
    this.hideLoadingTimer = setTimeout(() => this.resetValues(), 500);
  }

  resetValues() {
    this.loading = false;
    this.value = 0;
    this.align = 'start';
  }

  animateLoadingBar() {
    if(this.align == 'start') {
      if(this.value < 100)
        this.value = 100;
      else
        this.align = 'end';
    } else {
      if(this.value > 0)
        this.value = 0;
      else
        this.align = 'start';
    }
  }

}
