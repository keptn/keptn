import { Injectable, OnDestroy } from '@angular/core';
import { DataService } from './data.service';
import { KeptnInfo } from '../_models/keptn-info';
import { filter, takeUntil } from 'rxjs/operators';
import { IClientFeatureFlags } from '../../../shared/interfaces/feature-flags';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class FeatureFlagsService implements OnDestroy {
  public featureFlags?: IClientFeatureFlags;
  private unsubscribe$ = new Subject<void>();

  constructor(dataService: DataService) {
    dataService.keptnInfo
      .pipe(
        filter((info): info is KeptnInfo => !!info),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((info) => {
        this.featureFlags = info.bridgeInfo.featureFlags;
      });
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
