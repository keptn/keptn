import {Component, OnDestroy, OnInit} from '@angular/core';
import {Router, RoutesRecognized} from "@angular/router";
import {Observable, Subscription} from "rxjs";
import {filter, map} from "rxjs/operators";

import {Project} from "../_models/project";
import {DataService} from "../_services/data.service";

@Component({
  selector: 'app-header',
  templateUrl: './app-header.component.html',
  styleUrls: ['./app-header.component.scss']
})
export class AppHeaderComponent implements OnInit, OnDestroy {

  private routeSub: Subscription = Subscription.EMPTY;
  public projects: Observable<Project[]>;
  public project: Observable<Project>;

  constructor(private router: Router, private dataService: DataService) { }

  ngOnInit() {
    this.projects = this.dataService.projects;

    this.routeSub = this.router.events.subscribe(event => {
      if(event instanceof RoutesRecognized) {
        let projectName = event.state.root.children[0].params['projectName'];
        this.project = this.dataService.projects.pipe(
          filter(projects => !!projects),
          map(projects => projects.find(p => {
            return p.projectName === projectName;
          }))
        );
      }
    });
  }

  ngOnDestroy(): void {
    this.routeSub.unsubscribe();
  }

}
