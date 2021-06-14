import {ChangeDetectorRef, Component} from '@angular/core';
import {Observable} from "rxjs";

import {Project} from "../_models/project";

import {DataService} from "../_services/data.service";
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent {

  public projects$: Observable<Project[]>;

  public logoInvertedUrl = environment?.config?.logoInvertedUrl;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {
    this.projects$ = this.dataService.projects;
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

}
