import {Component, Input, OnDestroy, OnInit} from '@angular/core';

import {Trace} from '../../_models/trace';
import {Stage} from '../../_models/stage';

@Component({
  selector: 'ktb-stage-badge',
  templateUrl: './ktb-stage-badge.component.html',
  styleUrls: ['./ktb-stage-badge.component.scss'],
})
export class KtbStageBadgeComponent implements OnInit, OnDestroy {

  @Input() public evaluation: Trace;
  @Input() public stage: Stage;
  @Input() public isSelected: boolean | undefined = undefined;
  @Input() public fill: boolean | undefined = undefined;
  @Input() public error = false;
  @Input() public warning = false;
  @Input() public success = false;
  @Input() public highlight = false;

  constructor() { }

  ngOnInit() {
  }

  ngOnDestroy(): void {
  }

}
