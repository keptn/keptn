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

  getEventLabel(key: string): string {
    let label = key;
    switch(key) {
      case "sh.keptn.internal.event.service.create": {
        label = "Service create"
        break;
      }
      case "sh.keptn.event.configuration.change": {
        label = "Configuration change"
        break;
      }
      case "sh.keptn.event.monitoring.configure": {
        label = "Configure monitoring"
        break;
      }
      case "sh.keptn.events.deployment-finished": {
        label = "Deployment finished";
        break;
      }
      case "sh.keptn.events.tests-finished": {
        label = "Tests finished";
        break;
      }
      case "sh.keptn.events.evaluation-done": {
        label = "Evaluation done";
        break;
      }
      case "sh.keptn.internal.event.get-sli": {
        label = "Start SLI retrieval";
        break;
      }
      case "sh.keptn.internal.event.get-sli.done": {
        label = "SLI retrieval done";
        break;
      }
      case "sh.keptn.events.done": {
        label = "Done";
        break;
      }
      case "sh.keptn.events.done": {
        label = "Done";
        break;
      }


      case "sh.keptn.events.done": {
        label = "Done";
        break;
      }
      case "sh.keptn.event.problem.open": {
        label = "Problem open";
        break;
      }
      case "sh.keptn.events.problem": {
        label = "Problem detected";
        break;
      }
      default: {
        //statements;
        break;
      }
    }

    return label;
  }

  ngOnDestroy(): void {
    this._routeSubs.unsubscribe();
    this._projectSubs.unsubscribe();
  }

}
