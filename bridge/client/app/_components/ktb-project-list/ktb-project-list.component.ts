import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input } from '@angular/core';
import { Project } from '../../_models/project';

@Component({
  selector: 'ktb-project-list',
  templateUrl: './ktb-project-list.component.html',
  styleUrls: ['./ktb-project-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbProjectListComponent {

  public _projects: Project[] = [];

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
}
