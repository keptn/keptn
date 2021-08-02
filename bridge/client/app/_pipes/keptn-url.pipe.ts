import {Pipe, PipeTransform} from '@angular/core';
import {DataService} from '../_services/data.service';
import * as semver from 'semver';
import {Observable} from 'rxjs';
import {filter, take} from 'rxjs/operators';

@Pipe({
  name: 'keptnUrl'
})
export class KeptnUrlPipe implements PipeTransform {
  private static _version: Observable<any>;
  private static version: string;


  constructor(dataservice: DataService) {
    if (!KeptnUrlPipe._version) {
      KeptnUrlPipe._version = dataservice.keptnInfo;
      KeptnUrlPipe._version
        .pipe(
          filter(info => !!info && !!info.metadata),
          take(1)
        )
        .subscribe(info => {
          const vers = info.metadata.keptnversion;
          KeptnUrlPipe.version = `${semver.major(vers)}.${semver.minor(vers)}.x`;
        });
    }
  }

  transform(relativePath: string): string {
    return `https://keptn.sh/docs/${KeptnUrlPipe.version}${relativePath}`;
  }

}
