import { ChangeDetectorRef, Component, HostBinding } from '@angular/core';
import { filter, map, switchMap } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';

@Component({
  selector: 'ktb-environment-view',
  templateUrl: './ktb-environment-view.component.html',
  styleUrls: ['./ktb-environment-view.component.scss'],
  preserveWhitespaces: false,
})
export class KtbEnvironmentViewComponent {
  @HostBinding('class') cls = 'ktb-environment-view';
  public project$: Observable<Project | undefined>;

  constructor(
    private dataService: DataService,
    private route: ActivatedRoute,
    private _changeDetectorRef: ChangeDetectorRef
  ) {
    const projectName$ = this.route.paramMap.pipe(
      map((params) => params.get('projectName')),
      filter((projectName): projectName is string => !!projectName)
    );

    this.project$ = projectName$.pipe(
      switchMap((projectName) => this.dataService.getProject(projectName)),
      map((project) => (project?.isWholeProject ? project : undefined))
    );
  }
}
