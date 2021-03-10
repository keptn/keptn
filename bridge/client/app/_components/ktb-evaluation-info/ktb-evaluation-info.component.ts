import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {DtOverlayConfig} from "@dynatrace/barista-components/overlay";

import {Trace} from "../../_models/trace";

@Component({
  selector: 'ktb-evaluation-info',
  templateUrl: './ktb-evaluation-info.component.html',
  styleUrls: ['./ktb-evaluation-info.component.scss']
})
export class KtbEvaluationInfoComponent implements OnInit, OnDestroy {

  @Input() public evaluation: Trace;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };

  constructor() { }

  ngOnInit() {
  }

  ngOnDestroy(): void {
  }

}
