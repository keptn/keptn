import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { FeatureFlags } from '../../../shared/interfaces/feature-flags';

@Injectable({
  providedIn: 'root',
})
export class FeatureFlagsService {
  public featureFlags: FeatureFlags;

  constructor() {
    this.featureFlags = environment.featureFlags;
  }
}
