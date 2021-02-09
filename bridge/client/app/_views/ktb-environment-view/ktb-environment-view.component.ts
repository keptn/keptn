import {Component, Input, OnInit} from '@angular/core';
import {Project} from '../../_models/project';

@Component({
  selector: 'ktb-environment-view',
  templateUrl: './ktb-environment-view.component.html',
  styleUrls: ['./ktb-environment-view.component.scss']
})
export class KtbEnvironmentViewComponent implements OnInit {
  private _project: Project;

  @Input()
  get project() {
    return this._project;
  }
  set project(project: Project) {
    if (this._project !== project) {
      this._project = project;
    }
  }

  constructor() {
  }

  ngOnInit(): void {
  }
}
