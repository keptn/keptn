import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbHeatmapComponent } from './ktb-heatmap.component';
import { KtbHeatmapTooltipComponent } from './ktb-heatmap-tooltip.component';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbTestHeatmapComponent } from './testing/ktb-test-heatmap.component';

@NgModule({
  declarations: [KtbHeatmapComponent, KtbHeatmapTooltipComponent, KtbTestHeatmapComponent],
  imports: [CommonModule, DtKeyValueListModule, KtbPipeModule, DtIconModule, DtButtonModule],
  exports: [KtbHeatmapComponent, KtbHeatmapTooltipComponent, KtbTestHeatmapComponent],
})
export class KtbHeatmapModule {}
