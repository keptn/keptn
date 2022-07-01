import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { RootStoreFacade } from '../_stores/root/root.store.facade';
import { environment } from '../../environments/environment';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../_utils/app.utils';
import { Router } from '@angular/router';

@Component({
  selector: 'ktb-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  providers: [],
})
export class DashboardComponent implements OnInit, OnDestroy {
  refreshTimerSubscription = AppUtils.createTimer(0, this.initialDelayMillis).subscribe(() => this.refreshSequences());
  keptnMetadata$ = this.rootStoreFacade.metadata$;
  projects$ = this.rootStoreFacade.projects$;
  latestSequences$ = this.rootStoreFacade.latestSequences$;
  isQualityGatesOnly$ = this.rootStoreFacade.qualityGatesOnly$;
  logoInvertedUrl = environment?.config?.logoInvertedUrl;

  constructor(
    private router: Router,
    private rootStoreFacade: RootStoreFacade,
    @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number
  ) {
    const currentNav = this.router.getCurrentNavigation();
    const hadPreviousNavigation = currentNav != null && currentNav.previousNavigation != null;
    if (hadPreviousNavigation) {
      this.rootStoreFacade.refreshProjects();
    }
  }

  ngOnInit(): void {
    this.rootStoreFacade.refreshSequences();
  }

  refreshSequences(): void {
    this.rootStoreFacade.refreshSequences();
  }

  loadProjects(): void {
    this.rootStoreFacade.refreshProjects();
  }

  ngOnDestroy(): void {
    this.refreshTimerSubscription.unsubscribe();
  }
}
