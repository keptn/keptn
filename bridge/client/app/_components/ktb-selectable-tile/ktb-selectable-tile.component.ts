import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Directive,
  HostBinding,
  Input,
  ViewEncapsulation,
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
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSelectableTileComponent {
  @HostBinding('class') cls = 'ktb-selectable-tile';
  @HostBinding('attr.aria-selected') ariaSelected = this.selected;
  @HostBinding('class.ktb-tile-selected') tileSelected = this.selected;
  @HostBinding('attr.aria-disabled') ariaDisabled = this.disabled;
  @HostBinding('class.ktb-tile-disabled') tileDisabled = this.disabled;
  @HostBinding('attr.aria-error') ariaError = this.error;
  @HostBinding('class.ktb-tile-error') tileError = this.error;
  @HostBinding('attr.aria-warning') ariaWarning = this.warning;
  @HostBinding('class.ktb-tile-warning') tileWarning = this.warning;
  @HostBinding('attr.aria-success') ariaSuccess = this.success;
  @HostBinding('class.ktb-tile-success') tileSuccess = this.success;
  @HostBinding('attr.aria-highlight') ariaHighlight = this.highlight;
  @HostBinding('class.ktb-tile-highlight') tileHighlight = this.highlight;
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
