import { Component, OnDestroy } from '@angular/core';
import { filter, map, takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';

@Component({
  selector: 'ktb-no-service-info',
  templateUrl: './ktb-no-service-info.component.html',
  styleUrls: [],
})
export class KtbNoServiceInfoComponent implements OnDestroy {
  private unsubscribe$: Subject<void> = new Subject();
  public createServiceLink = '';

  constructor(router: ActivatedRoute, public readonly location: Location) {
    router.paramMap
      .pipe(
        map((params) => params.get('projectName')),
        filter((projectName): projectName is string => !!projectName),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((projectName) => {
        this.createServiceLink = `/project/${projectName}/settings/services/create`;
      });
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
