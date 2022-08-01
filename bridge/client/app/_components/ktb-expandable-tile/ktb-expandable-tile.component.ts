import {
  ChangeDetectionStrategy,
  Component,
  Directive,
  EventEmitter,
  HostBinding,
  Input,
  Output,
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

  @Output()
  public expandedChange = new EventEmitter<boolean>();

  @Input()
  alignment: alignmentType = 'right';

  @Input()
  @HostBinding('attr.aria-error')
  @HostBinding('class.ktb-tile-error')
  error = false;

  @Input()
  @HostBinding('attr.aria-success')
  @HostBinding('class.ktb-tile-success')
  success = false;

  @Input()
  expanded = false;

  /** Whether the tile is disabled. */
  @Input()
  @HostBinding('attr.aria-disabled')
  @HostBinding('class.ktb-tile-disabled')
  disabled = false;

  @Input()
  @HostBinding('attr.aria-warning')
  @HostBinding('class.ktb-tile-warning')
  warning = false;

  @Input()
  @HostBinding('attr.aria-highlight')
  @HostBinding('class.ktb-tile-highlight')
  highlight = false;

  toggle(): void {
    this.expanded = !this.expanded;
    this.expandedChange.emit(this.expanded);
  }
}
