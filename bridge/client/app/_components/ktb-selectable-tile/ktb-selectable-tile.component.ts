import { Component, Directive, HostBinding, Input, ViewEncapsulation } from '@angular/core';

@Directive({
  selector: `ktb-selectable-tile-header, [ktb-selectable-tile-header], [ktbSelectableTileHeader]`,
  exportAs: 'ktbSelectableTileHeader',
})
export class KtbSelectableTileHeaderDirective {}

@Component({
  selector: 'ktb-selectable-tile',
  templateUrl: './ktb-selectable-tile.component.html',
  styleUrls: ['./ktb-selectable-tile.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
})
export class KtbSelectableTileComponent {
  @HostBinding('class') cls = 'ktb-selectable-tile';
  private _selected = false;
  private _disabled = false;
  private _error = false;
  private _warning = false;
  private _success = false;
  private _highlight = false;

  /** Whether the tile is selected. */
  @Input()
  @HostBinding('class.ktb-tile-selected')
  @HostBinding('class.ktb-tile-selected')
  get selected(): boolean {
    return this._selected && !this.disabled;
  }

  set selected(value: boolean) {
    if (this._selected !== value) {
      this._selected = value;
    }
  }

  /** Whether the tile is disabled. */
  @Input()
  @HostBinding('attr.aria-disabled')
  @HostBinding('class.ktb-tile-disabled')
  get disabled(): boolean {
    return this._disabled;
  }

  set disabled(value: boolean) {
    if (this._disabled && this._selected) {
      this._selected = false;
    }
  }

  @Input()
  @HostBinding('class.ktb-tile-error')
  @HostBinding('attr.aria-error')
  get error(): boolean {
    return this._error;
  }

  set error(value: boolean) {
    if (this._error !== value) {
      this._error = value;
    }
  }

  @Input()
  @HostBinding('attr.aria-warning')
  @HostBinding('class.ktb-tile-warning')
  get warning(): boolean {
    return this._warning;
  }

  set warning(value: boolean) {
    if (this._warning !== value) {
      this._warning = value;
    }
  }

  @Input()
  @HostBinding('attr.aria-success')
  @HostBinding('class.ktb-tile-success')
  get success(): boolean {
    return this._success;
  }

  set success(value: boolean) {
    if (this._success !== value) {
      this._success = value;
    }
  }

  @Input()
  @HostBinding('attr.aria-highlight')
  @HostBinding('class.ktb-tile-highlight')
  get highlight(): boolean {
    return this._highlight;
  }

  set highlight(value: boolean) {
    if (this._highlight !== value) {
      this._highlight = value;
    }
  }
}
