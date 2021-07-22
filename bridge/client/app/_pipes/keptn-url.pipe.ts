import {Pipe, PipeTransform} from '@angular/core';
import {DataService} from '../_services/data.service';
import semver from 'semver';
import {Observable} from 'rxjs';
import {filter, take} from 'rxjs/operators';
import { KeptnInfo } from '../_models/keptn-info';

@Pipe({
  name: 'keptnUrl'
})
export class KeptnUrlPipe implements PipeTransform {
  private static _version: Observable<KeptnInfo | undefined>;
  private static version: string;


  constructor(dataService: DataService) {
    if (!KeptnUrlPipe._version) {
      KeptnUrlPipe._version = dataService.keptnInfo;
      KeptnUrlPipe._version
        .pipe(
          filter((info: KeptnInfo | undefined): info is KeptnInfo => !!info?.metadata),
          take(1)
        )
        .subscribe(info => {
          const version = info.metadata.keptnversion;
          KeptnUrlPipe.version = `${semver.major(version)}.${semver.minor(version)}.x`;
        });
    }
  }

  transform(relativePath: string): string {
    return `https://keptn.sh/docs/${KeptnUrlPipe.version}${relativePath}`;
  }

}
