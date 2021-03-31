import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'ktb-sli-breakdown-criteria-item',
  templateUrl: './ktb-sli-breakdown-criteria-item.component.html',
  styleUrls: ['./ktb-sli-breakdown-criteria-item.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbSliBreakdownCriteriaItemComponent implements OnInit {
  private _targets: any;

  @Input()
  get targets() {
    return this._targets;
  }

  set targets(targets: any) {
    if(this._targets !== targets) {
      this._targets = targets;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit(): void {
  }

}
