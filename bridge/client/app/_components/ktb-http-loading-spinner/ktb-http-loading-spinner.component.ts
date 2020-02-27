import {Component, Input, OnInit} from '@angular/core';
import {HttpStateService} from "../../_services/http-state.service";
import {HttpProgressState, HttpState} from "../../_models/http-progress-state";

@Component({
  selector: 'ktb-http-loading-spinner',
  templateUrl: './ktb-http-loading-spinner.component.html',
  styleUrls: ['./ktb-http-loading-spinner.component.scss']
})
export class KtbHttpLoadingSpinnerComponent implements OnInit {

  public loading = false;
  @Input() public filterBy: string | null = null;

  private hideLoadingTimer;

  constructor(private httpStateService: HttpStateService) { }

  ngOnInit() {
    this.httpStateService.state.subscribe((progress: HttpState) => {
      if (progress && progress.url) {
        if(!this.filterBy || progress.url.indexOf(this.filterBy) !== -1) {
          if(progress.state === HttpProgressState.start)
            this.showLoadingSpinner();
          else
            this.hideLoadingSpinner();
        }
      }
    });
  }

  showLoadingSpinner() {
    clearTimeout(this.hideLoadingTimer);
    this.loading = true;
  }

  hideLoadingSpinner() {
    this.hideLoadingTimer = setTimeout(() => { this.loading = false; }, 500);
  }

}
