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

  public projects$: Observable<Project[]>;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {
    this.projects$ = this.dataService.projects;
  }

  ngOnInit() {
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  ngOnDestroy(): void {
  }

}
