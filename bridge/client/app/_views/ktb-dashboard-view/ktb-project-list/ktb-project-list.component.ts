import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { Sequence } from '../../../_models/sequence';
import { IMetadata } from '../../../_interfaces/metadata';
import { IProject } from '../../../../../shared/interfaces/project';

export type ProjectSequences = Record<string, Sequence[]>;

@Component({
  selector: 'ktb-project-list',
  templateUrl: './ktb-project-list.component.html',
  styleUrls: ['./ktb-project-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbProjectListComponent {
  @Input() metadata: IMetadata | null = null;
  @Input() projects: IProject[] = [];
  @Input() sequences: ProjectSequences = {};

  getSequencesPerProject(project: IProject): Sequence[] {
    const latestSequences = this.sequences[project.projectName];
    return latestSequences ?? [];
  }
}
