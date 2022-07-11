import { Component } from '@angular/core';
import { mergeMap, Observable } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { filter, finalize, map } from 'rxjs/operators';
import { DataService } from '../../../_services/data.service';

@Component({
  selector: 'ktb-service-settings-overview',
  templateUrl: './ktb-service-settings-overview.component.html',
})
export class KtbServiceSettingsOverviewComponent {
  isLoading = true;
  projectName$: Observable<string> = this.route.paramMap.pipe(
    map((params) => params.get('projectName')),
    filter((projectName): projectName is string => !!projectName)
  );
  serviceNames$: Observable<string[]> = this.projectName$.pipe(
    mergeMap((projectName) => this.dataService.getServiceNames(projectName)),
    finalize(() => (this.isLoading = false))
  );

  constructor(private route: ActivatedRoute, private dataService: DataService) {}
}
