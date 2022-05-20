import { ChangeDetectionStrategy, ChangeDetectorRef, Component, ElementRef, Input } from '@angular/core';
import { AppUtils } from '../../_utils/app.utils';
import { IHeatmapTooltip, IHeatmapTooltipType } from '../../_interfaces/heatmap';

@Component({
  selector: 'ktb-heatmap-tooltip',
  templateUrl: './ktb-heatmap-tooltip.component.html',
  styleUrls: ['./ktb-heatmap-tooltip.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbHeatmapTooltipComponent {
  public readonly IHeatmapTooltipType = IHeatmapTooltipType;
  public readonly formatNumber = AppUtils.formatNumber;
  private _tooltip?: IHeatmapTooltip;

  @Input()
  set tooltip(tooltip: IHeatmapTooltip | undefined) {
    if (this._tooltip !== tooltip) {
      this._tooltip = tooltip;
      this._changeDetectorRef.detectChanges();
    }
  }
  get tooltip(): IHeatmapTooltip | undefined {
    return this._tooltip;
  }

  constructor(public _elementRef: ElementRef, private _changeDetectorRef: ChangeDetectorRef) {}
}
