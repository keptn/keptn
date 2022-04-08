import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'ktb-loading-spinner',
  templateUrl: './ktb-loading-spinner.component.html',
  styleUrls: ['./ktb-loading-spinner.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbLoadingSpinnerComponent {}
