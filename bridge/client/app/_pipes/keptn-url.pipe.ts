import { Pipe, PipeTransform } from '@angular/core';
import { DataService } from '../_services/data.service';
import semver from 'semver';
import { Observable } from 'rxjs';
import { filter, take } from 'rxjs/operators';
import { KeptnInfo } from '../_models/keptn-info';

@Pipe({
  name: 'keptnUrl',
})
export class KeptnUrlPipe implements PipeTransform {
  private static _version?: Observable<KeptnInfo | undefined>;
  private static version = '';

  constructor(dataService: DataService) {
    if (!KeptnUrlPipe._version) {
      KeptnUrlPipe._version = dataService.keptnInfo;
      KeptnUrlPipe._version
        .pipe(
          filter((info: KeptnInfo | undefined): info is KeptnInfo => !!info),
          take(1)
        )
        .subscribe((info) => {
          const version = info.bridgeInfo.bridgeVersion;
          if (!version) {
            return;
          }
          KeptnUrlPipe.version = semver.valid(version) ? `${semver.major(version)}.${semver.minor(version)}.x` : '';
        });
    }
  }

  transform(relativePath: string): string {
    return KeptnUrlPipe.version
      ? `https://keptn.sh/docs/${KeptnUrlPipe.version}${relativePath}`
      : 'https://keptn.sh/docs/install/';
  }
}
