import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input} from '@angular/core';

@Component({
  selector: 'ktb-sli-breakdown-criteria-item',
  templateUrl: './ktb-sli-breakdown-criteria-item.component.html',
  styleUrls: [],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbSliBreakdownCriteriaItemComponent {
  private _targets: any;

  @Input()
  get targets() {
    return this._targets;
  }

  @Input()
  public isInformative = false;

  set targets(targets: any) {
    if(this._targets !== targets) {
      this._targets = targets;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) {}
}
