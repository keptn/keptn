import { Injectable } from '@angular/core';
import { DataService } from './data.service';
import { KeptnInfo } from '../_models/keptn-info';
import { filter } from 'rxjs/operators';
import { IClientFeatureFlags } from '../../../shared/interfaces/feature-flags';

@Injectable({
  providedIn: 'root',
})
export class FeatureFlagsService {
  public featureFlags?: IClientFeatureFlags;

  constructor(dataService: DataService) {
    dataService.keptnInfo.pipe(filter((info): info is KeptnInfo => !!info)).subscribe((info) => {
      this.featureFlags = info.bridgeInfo.featureFlags;
    });
  }
}
