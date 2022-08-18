import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { SequenceState } from '../../../_models/sequenceState';
import { IMetadata } from '../../../_interfaces/metadata';
import { IProject } from '../../../../../shared/interfaces/project';

export type ProjectSequences = Record<string, SequenceState[]>;

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

  getSequencesPerProject(project: IProject): SequenceState[] {
    const latestSequences = this.sequences[project.projectName];
    return latestSequences ?? [];
  }
}
