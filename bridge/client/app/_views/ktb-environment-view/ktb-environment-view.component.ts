import {Component, OnDestroy, OnInit} from '@angular/core';
import {map, switchMap} from "rxjs/operators";
import {Observable, Subject} from "rxjs";
import {ActivatedRoute} from "@angular/router";

import {Project} from '../../_models/project';
import {DataService} from "../../_services/data.service";

@Component({
  selector: 'ktb-environment-view',
  templateUrl: './ktb-environment-view.component.html',
  styleUrls: ['./ktb-environment-view.component.scss'],
  host: {
    class: 'ktb-environment-view'
  },
  preserveWhitespaces: false
})
export class KtbEnvironmentViewComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  public project$: Observable<Project>;

  constructor(private dataService: DataService, private route: ActivatedRoute) {
  }

  ngOnInit(): void {
    this.project$ = this.route.params
      .pipe(
        map(params => params.projectName),
        switchMap(projectName => this.dataService.getProject(projectName))
      );
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
