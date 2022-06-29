import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { IProject } from '../../../../shared/interfaces/project';
import { getDistinctServiceNames, getShipyardVersion, isShipyardNotSupported } from '../../_models/project';
import { ISequence } from '../../../../shared/interfaces/sequence';

@Component({
  selector: 'ktb-project-tile',
  templateUrl: './ktb-project-tile.component.html',
  styleUrls: ['./ktb-project-tile.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbProjectTileComponent {
  @Input() project?: IProject;
  @Input() sequences: ISequence[] = [];
  @Input() supportedShipyardVersion: string | undefined;

  getShipyardVersion = getShipyardVersion;
  isShipyardNotSupported = isShipyardNotSupported;
  getDistinctServiceNames = getDistinctServiceNames;
}
