import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnInit, ViewEncapsulation} from '@angular/core';
import {coerceBooleanProperty} from "@angular/cdk/coercion";

@Component({
  selector: 'ktb-selectable-tile',
  templateUrl: './ktb-selectable-tile.component.html',
  styleUrls: ['./ktb-selectable-tile.component.scss'],
  host: {
    class: 'ktb-selectable-tile',
    '[attr.aria-selected]': 'selected',
    '[class.ktb-tile-selected]': 'selected',
    '[attr.aria-disabled]': 'disabled',
    '[class.ktb-tile-disabled]': 'disabled',
    '[attr.aria-error]': 'error',
    '[class.ktb-tile-error]': 'error',
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSelectableTileComponent implements OnInit {

  private _selected = false;
  private _disabled = false;
  private _error = false;

  /** Whether the tile is selected. */
  @Input()
  get selected(): boolean {
    return this._selected && !this.disabled;
  }
  set selected(value: boolean) {
    const newValue = coerceBooleanProperty(value);
    if (this._selected !== newValue) {
      this._selected = newValue;
      this._changeDetectorRef.markForCheck();
    }
  }

  /** Whether the tile is disabled. */
  @Input()
  get disabled(): boolean {
    return this._disabled;
  }
  set disabled(value: boolean) {
    this._disabled = coerceBooleanProperty(value);
    if (this._disabled && this._selected) {
      this._selected = false;
      this._changeDetectorRef.markForCheck();
    }
  }

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

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }

}
