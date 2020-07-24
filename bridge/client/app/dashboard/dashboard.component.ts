import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {Observable, Subject, Subscription} from "rxjs";

import {Project} from "../_models/project";

import {DataService} from "../_services/data.service";
import {takeUntil} from "rxjs/operators";

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();

  public projects$: Observable<Project[]>;
  public error: boolean = false;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {
  }

  ngOnInit() {
    this.projects$ = this.dataService.projects;
    this.projects$
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(projects => {
        this.error = false;
      }, error => {
        this.error = true;
      });
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
