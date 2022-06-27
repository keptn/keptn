import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { Project } from '../../_models/project';
import { Sequence } from '../../_models/sequence';

@Component({
  selector: 'ktb-project-tile',
  templateUrl: './ktb-project-tile.component.html',
  styleUrls: ['./ktb-project-tile.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbProjectTileComponent {
  @Input() project?: Project;
  @Input() sequences: Sequence[] = [];
  @Input() supportedShipyardVersion: string | undefined;
}
