import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Directive,
  HostBinding,
  Input,
  ViewEncapsulation,
} from '@angular/core';

type alignmentType = 'right' | 'left';
@Directive({
  selector: `ktb-expandable-tile-header, [ktb-expandable-tile-header], [ktbExpandableTileHeader]`,
  exportAs: 'ktbExpandableTileHeader',
})
export class KtbExpandableTileHeaderDirective {}

@Component({
  selector: 'ktb-expandable-tile',
  templateUrl: './ktb-expandable-tile.component.html',
  styleUrls: ['./ktb-expandable-tile.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbExpandableTileComponent {
  @HostBinding('class') cls = 'ktb-expandable-tile';
  private _error = false;
  private _success = false;
  private _expanded = false;
  private _disabled = false;
  private _warning = false;
  private _highlight = false;
  private _alignment: alignmentType = 'right';

  @Input()
  get alignment(): alignmentType {
    return this._alignment;
  }
  set alignment(alignment: alignmentType) {
    this._alignment = alignment;
  }

  @Input()
  @HostBinding('attr.aria-error')
  @HostBinding('class.ktb-tile-error')
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
  @HostBinding('attr.aria-success')
  @HostBinding('class.ktb-tile-success')
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
  get expanded(): boolean {
    return this._expanded;
  }
  set expanded(value: boolean) {
    if (this._expanded !== value) {
      this._expanded = value;
      this._changeDetectorRef.markForCheck();
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
    if (this._disabled !== value) {
      this._disabled = value;
      this._changeDetectorRef.markForCheck();
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
      this._changeDetectorRef.markForCheck();
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
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) {}

  toggle(): void {
    this.expanded = !this.expanded;
  }
}
