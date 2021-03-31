import {ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {takeUntil} from "rxjs/operators";
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
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbEnvironmentViewComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  public project$: Observable<Project>;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute) {
  }

  ngOnInit(): void {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.project$ = this.dataService.getProject(params.projectName);
      });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
