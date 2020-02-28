import {Component, OnDestroy, OnInit} from '@angular/core';
import {Observable, Subscription} from "rxjs";

import {Project} from "../_models/project";

import {DataService} from "../_services/data.service";

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, OnDestroy {

  public projects: Observable<Project[]>;
  public error: boolean = false;

  private _projectsSubs: Subscription = Subscription.EMPTY;

  constructor(private dataService: DataService) {
  }

  ngOnInit() {
    this.projects = this.dataService.projects;
    this._projectsSubs = this.projects.subscribe(projects => {
      this.error = false;
    }, error => {
      this.error = true;
    });
  }

  ngOnDestroy(): void {
    this._projectsSubs.unsubscribe();
  }

}
