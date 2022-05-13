import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbHeatmapComponent } from './ktb-heatmap.component';
import { KtbHeatmapTooltipComponent } from './ktb-heatmap-tooltip.component';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbTestHeatmapComponent } from './testing/ktb-test-heatmap.component';
import { HttpClientModule } from '@angular/common/http';

@NgModule({
  declarations: [KtbHeatmapComponent, KtbHeatmapTooltipComponent, KtbTestHeatmapComponent], // add KtbTestHeatmapComponent for testing
  imports: [
    CommonModule,
    DtKeyValueListModule,
    KtbPipeModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    HttpClientModule, // for icons
    DtButtonModule,
  ],
  exports: [KtbHeatmapComponent, KtbHeatmapTooltipComponent], // add KtbTestHeatmapComponent for testing
})
export class KtbHeatmapModule {}
