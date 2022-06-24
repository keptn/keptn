import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { Project } from '../../_models/project';
import { Sequence } from '../../_models/sequence';
import { IMetadata } from '../../_interfaces/metadata';

export type ProjectSequences = Record<string, Sequence[]>;

@Component({
  selector: 'ktb-project-list',
  templateUrl: './ktb-project-list.component.html',
  styleUrls: ['./ktb-project-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbProjectListComponent {
  @Input() metadata: IMetadata | null = null;
  @Input() projects: Project[] = [];
  @Input() sequences: ProjectSequences = {};

  getSequencesPerProject(project: Project): Sequence[] {
    const latestSequences = this.sequences[project.projectName];
    return latestSequences ?? [];
  }

  getShipyardversion(): string | undefined {
    return this.metadata != null ? this.metadata?.shipyardversion : undefined;
  }
}
