import { Component, HostBinding } from '@angular/core';
import { map, switchMap } from 'rxjs/operators';
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

  constructor(private dataService: DataService, private route: ActivatedRoute) {
    this.project$ = this.route.params.pipe(
      map((params) => params.projectName),
      switchMap((projectName) => this.dataService.getProject(projectName))
    );
  }
}
