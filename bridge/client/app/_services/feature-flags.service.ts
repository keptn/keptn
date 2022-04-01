import { Injectable, OnDestroy } from '@angular/core';
import { DataService } from './data.service';
import { KeptnInfo } from '../_models/keptn-info';
import { filter, map, takeUntil } from 'rxjs/operators';
import { IClientFeatureFlags } from '../../../shared/interfaces/feature-flags';
import { Observable, Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class FeatureFlagsService implements OnDestroy {
  public featureFlags$: Observable<IClientFeatureFlags>;
  private unsubscribe$ = new Subject<void>();

  constructor(dataService: DataService) {
    this.featureFlags$ = dataService.keptnInfo.pipe(
      filter((info): info is KeptnInfo => !!info),
      map((info) => info.bridgeInfo.featureFlags),
      takeUntil(this.unsubscribe$)
    );
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
