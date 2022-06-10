import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

@Component({
  selector: 'ktb-loading-spinner',
  templateUrl: './ktb-loading-spinner.component.html',
  styleUrls: ['./ktb-loading-spinner.component.scss'],
  host: {
    role: 'progressbar',
    'aria-busy': 'true',
    'aria-live': 'assertive',
    '[attr.aria-label]': 'ariaLabel',
  },
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbLoadingSpinnerComponent {
  /** The aria-label attribute. */
  @Input('aria-label') ariaLabel?: string;
}
