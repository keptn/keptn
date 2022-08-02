import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { SequenceState } from '../../../_models/sequenceState';
import { IProject } from '../../../../../shared/interfaces/project';
import { getDistinctServiceNames, getShipyardVersion, isShipyardNotSupported } from '../../../_models/project';

@Component({
  selector: 'ktb-project-tile',
  templateUrl: './ktb-project-tile.component.html',
  styleUrls: ['./ktb-project-tile.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbProjectTileComponent {
  @Input() project?: IProject;
  @Input() sequences: SequenceState[] = [];
  @Input() supportedShipyardVersion: string | undefined;

  getShipyardVersion = getShipyardVersion;
  isShipyardNotSupported = isShipyardNotSupported;
  getDistinctServiceNames = getDistinctServiceNames;
}
