import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbHeatmapTooltipComponent } from './ktb-heatmap-tooltip.component';
import { KtbHeatmapComponent } from './ktb-heatmap.component';

@NgModule({
  declarations: [KtbHeatmapComponent, KtbHeatmapTooltipComponent],
  imports: [
    CommonModule,
    HttpClientModule,
    DtKeyValueListModule,
    KtbPipeModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtButtonModule,
  ],
  exports: [KtbHeatmapComponent, KtbHeatmapTooltipComponent],
})
export class KtbHeatmapModule {}
