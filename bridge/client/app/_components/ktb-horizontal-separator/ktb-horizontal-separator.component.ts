import { ChangeDetectionStrategy, Component, Directive, HostBinding, ViewEncapsulation } from '@angular/core';

@Directive({
  selector: `ktb-horizontal-separator-title, [ktb-horizontal-separator-title], [ktbHorizontalSeparatorTitle]`,
  exportAs: 'ktbHorizontalSeparatorTitle',
})
export class KtbHorizontalSeparatorTitleDirective {
  @HostBinding('class') cls = 'ktb-horizontal-separator-title';
}

@Component({
  selector: 'ktb-horizontal-separator',
  templateUrl: './ktb-horizontal-separator.component.html',
  styleUrls: ['./ktb-horizontal-separator.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbHorizontalSeparatorComponent {}
