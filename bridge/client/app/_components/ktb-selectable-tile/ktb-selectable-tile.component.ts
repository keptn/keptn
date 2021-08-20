import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Directive,
  Input,
  ViewEncapsulation
} from '@angular/core';

@Directive({
  selector: `ktb-selectable-tile-header, [ktb-selectable-tile-header], [ktbSelectableTileHeader]`,
  exportAs: 'ktbSelectableTileHeader',
})
export class KtbSelectableTileHeaderDirective {}

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
    '[attr.aria-warning]': 'warning',
    '[class.ktb-tile-warning]': 'warning',
    '[attr.aria-success]': 'success',
    '[class.ktb-tile-success]': 'success',
    '[attr.aria-highlight]': 'highlight',
    '[class.ktb-tile-highlight]': 'highlight',
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSelectableTileComponent {

  private _selected = false;
  private _disabled = false;
  private _error = false;
  private _warning = false;
  private _success = false;
  private _highlight = false;

  /** Whether the tile is selected. */
  @Input()
  get selected(): boolean {
    return this._selected && !this.disabled;
  }
  set selected(value: boolean) {
    if (this._selected !== value) {
      this._selected = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  /** Whether the tile is disabled. */
  @Input()
  get disabled(): boolean {
    return this._disabled;
  }
  set disabled(value: boolean) {
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
    if (this._error !== value) {
      this._error = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get warning(): boolean {
    return this._warning;
  }
  set warning(value: boolean) {
    if (this._warning !== value) {
      this._warning = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get success(): boolean {
    return this._success;
  }
  set success(value: boolean) {
    if (this._success !== value) {
      this._success = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get highlight(): boolean {
    return this._highlight;
  }
  set highlight(value: boolean) {
    if (this._highlight !== value) {
      this._highlight = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) {}
}
