import { Component, Input } from '@angular/core';
import { Project } from '../../_models/project';
import { Sequence } from '../../_models/sequence';

@Component({
  selector: 'ktb-project-tile',
  templateUrl: './ktb-project-tile.component.html',
  styleUrls: ['./ktb-project-tile.component.scss'],
})
export class KtbProjectTileComponent {
  private _project?: Project;
  private _sequences: Sequence[] = [];
  private _supportedShipyardVersion: string | undefined;

  @Input()
  get project(): Project | undefined {
    return this._project;
  }

  set project(value: Project | undefined) {
    if (this._project !== value) {
      this._project = value;
    }
  }

  @Input()
  get sequences(): Sequence[] {
    return this._sequences;
  }

  set sequences(value: Sequence[]) {
    this._sequences = value;
  }

  @Input()
  get supportedShipyardVersion(): string | undefined {
    return this._supportedShipyardVersion;
  }

  set supportedShipyardVersion(value: string | undefined) {
    this._supportedShipyardVersion = value;
  }
}
