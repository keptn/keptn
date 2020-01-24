import {Component, OnDestroy, OnInit} from '@angular/core';
import {first, map} from "rxjs/operators";
import {Observable, Subscription} from "rxjs";
import {ActivatedRoute} from "@angular/router";

import {Root} from "../_models/root";
import {Project} from "../_models/project";

import {DataService} from "../_services/data.service";

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {

  public project: Observable<Project>;
  public currentRoot: Root;
  public error: boolean = false;

  private _routeSubs: Subscription = Subscription.EMPTY;
  private _projectSubs: Subscription = Subscription.EMPTY;

  constructor(private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit() {
    this._routeSubs = this.route.params.subscribe(params => {
      if(params['projectName']) {
        this.currentRoot = null;

        this.project = this.dataService.projects.pipe(
          map(projects => projects.find(project => {
            return project.projectName === params['projectName'];
          }))
        );

        this.project
          .pipe(first(project => !!project && !!project.getServices()))
          .subscribe(project => {
            project.getServices().forEach(service => {
              this.dataService.loadRoots(project, service);
            });
          }, error => {
            this.error = true;
          });
      }
    });
  }

  loadTraces(root): void {
    this.dataService.loadTraces(root);
  }

  ngOnDestroy(): void {
    this._routeSubs.unsubscribe();
    this._projectSubs.unsubscribe();
  }

}
