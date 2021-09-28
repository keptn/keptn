import { Component } from '@angular/core';
import { map, switchMap } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';

@Component({
  selector: 'ktb-environment-view',
  templateUrl: './ktb-environment-view.component.html',
  styleUrls: ['./ktb-environment-view.component.scss'],
  host: {
    class: 'ktb-environment-view',
  },
  preserveWhitespaces: false,
})
export class KtbEnvironmentViewComponent {
  public project$: Observable<Project | undefined>;

  constructor(private dataService: DataService, private route: ActivatedRoute) {
    this.project$ = this.route.params.pipe(
      map((params) => params.projectName),
      switchMap((projectName) => this.dataService.getProject(projectName))
    );
  }
}
