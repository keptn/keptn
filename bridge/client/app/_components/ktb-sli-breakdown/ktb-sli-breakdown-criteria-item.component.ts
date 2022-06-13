import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input } from '@angular/core';
import { Target } from '../../../../shared/interfaces/indicator-result';

@Component({
  selector: 'ktb-sli-breakdown-criteria-item',
  templateUrl: './ktb-sli-breakdown-criteria-item.component.html',
  styleUrls: [],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSliBreakdownCriteriaItemComponent {
  private _targets: Target[] = [];

  @Input()
  get targets(): Target[] {
    return this._targets;
  }
  set targets(targets: Target[]) {
    if (this._targets !== targets) {
      this._targets = targets;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  public isInformative = false;

  constructor(private _changeDetectorRef: ChangeDetectorRef) {}
}
