import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Directive,
  Input,
  OnInit,
  ViewEncapsulation
} from '@angular/core';
import {coerceBooleanProperty} from "@angular/cdk/coercion";

@Directive({
  selector: `ktb-expandable-tile-header, [ktb-expandable-tile-header], [ktbExpandableTileHeader]`,
  exportAs: 'ktbExpandableTileHeader',
})
export class KtbExpandableTileHeader {}

@Component({
  selector: 'ktb-expandable-tile',
  templateUrl: './ktb-expandable-tile.component.html',
  styleUrls: ['./ktb-expandable-tile.component.scss'],
  host: {
    class: 'ktb-expandable-tile',
    '[attr.aria-error]': 'error',
    '[class.ktb-tile-error]': 'error',
    '[attr.aria-success]': 'success',
    '[class.ktb-tile-success]': 'success',
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbExpandableTileComponent implements OnInit {

  private _error = false;
  private _success = false;

  @Input()
  get error(): boolean {
    return this._error;
  }
  set error(value: boolean) {
    const newValue = coerceBooleanProperty(value);
    if (this._error !== newValue) {
      this._error = newValue;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get success(): boolean {
    return this._success;
  }
  set success(value: boolean) {
    const newValue = coerceBooleanProperty(value);
    if (this._success !== newValue) {
      this._success = newValue;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }

}
