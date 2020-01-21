import {ChangeDetectorRef, Component, Input, OnInit} from '@angular/core';
import {Project} from "../../_models/project";

@Component({
  selector: 'ktb-project-tile',
  templateUrl: './ktb-project-tile.component.html',
  styleUrls: ['./ktb-project-tile.component.scss']
})
export class KtbProjectTileComponent implements OnInit {

  public _project: Project;

  @Input()
  get project(): Project {
    return this._project;
  }
  set project(value: Project) {
    if (this._project !== value) {
      this._project = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }
}
