import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';

@Component({
  selector: 'ktb-service-settings-overview',
  templateUrl: './ktb-service-settings-overview.component.html',
})
export class KtbServiceSettingsOverviewComponent {
  public projectName$: Observable<string | null>;

  constructor(route: ActivatedRoute) {
    this.projectName$ = route.paramMap.pipe(
      map(params => params.get('projectName')),
    );
  }

}
