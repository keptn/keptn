import {ChangeDetectorRef, Component, Input, OnInit} from '@angular/core';
import {Project} from "../../_models/project";

@Component({
  selector: 'ktb-project-list',
  templateUrl: './ktb-project-list.component.html',
  styleUrls: ['./ktb-project-list.component.scss']
})
export class KtbProjectListComponent implements OnInit {

  public _projects: Project[];

  @Input()
  get projects(): Project[] {
    return this._projects;
  }
  set projects(value: Project[]) {
    if (this._projects !== value) {
      this._projects = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }
}
