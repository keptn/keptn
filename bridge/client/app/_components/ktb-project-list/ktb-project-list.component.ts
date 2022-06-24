import { Component, Input } from '@angular/core';
import { Project } from '../../_models/project';
import { Sequence } from '../../_models/sequence';
import { IMetadata } from '../../_interfaces/metadata';

export type ProjectSequences = Record<string, Sequence[]>;

@Component({
  selector: 'ktb-project-list',
  templateUrl: './ktb-project-list.component.html',
  styleUrls: ['./ktb-project-list.component.scss'],
})
export class KtbProjectListComponent {
  private _metadata: IMetadata | null = null;
  private _projects: Project[] = [];
  private _sequences: ProjectSequences = {};

  @Input()
  get projects(): Project[] {
    return this._projects;
  }

  set projects(value: Project[]) {
    if (this._projects !== value) {
      this._projects = value;
    }
  }

  @Input()
  get metadata(): IMetadata | null {
    return this._metadata;
  }

  set metadata(value: IMetadata | null) {
    this._metadata = value;
  }

  @Input()
  get sequences(): ProjectSequences {
    return this._sequences;
  }

  set sequences(value: ProjectSequences) {
    this._sequences = value;
  }

  getSequencesPerProject(project: Project): Sequence[] {
    const latestSequences = this._sequences[project.projectName];
    return latestSequences ?? [];
  }

  getShipyardversion(): string | undefined {
    return this.metadata != null ? this._metadata?.shipyardversion : undefined;
  }
}
