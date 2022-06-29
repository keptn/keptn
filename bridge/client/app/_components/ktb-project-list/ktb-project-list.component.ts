import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { IMetadata } from '../../_interfaces/metadata';
import { IProject } from '../../../../shared/interfaces/project';
import { ISequence } from '../../../../shared/interfaces/sequence';

export type ProjectSequences = Record<string, ISequence[]>;

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

  getSequencesPerProject(project: IProject): ISequence[] {
    const latestSequences = this.sequences[project.projectName];
    return latestSequences ?? [];
  }
}
