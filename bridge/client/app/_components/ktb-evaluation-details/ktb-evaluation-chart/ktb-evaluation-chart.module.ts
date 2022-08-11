import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtChartModule } from '@dynatrace/barista-components/chart';
import { KtbEvaluationChartComponent } from './ktb-evaluation-chart.component';
import { KtbEvaluationChartLegacyComponent } from './ktb-evaluation-chart-legacy/ktb-evaluation-chart-legacy.component';
import { DateFormatPipe, MomentModule } from 'ngx-moment';
import { KtbHeatmapModule } from '../../ktb-heatmap/ktb-heatmap.module';
import { KtbChartModule } from '../../ktb-chart/ktb-chart.module';
import { DtButtonGroupModule } from '@dynatrace/barista-components/button-group';
import { FlexModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';

@NgModule({
  declarations: [KtbEvaluationChartComponent, KtbEvaluationChartLegacyComponent],
  imports: [
    CommonModule,
    DtButtonGroupModule,
    DtButtonModule,
    DtChartModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtKeyValueListModule,
    FlexModule,
    KtbChartModule,
    KtbHeatmapModule,
    KtbPipeModule,
    MomentModule,
  ],
  exports: [KtbEvaluationChartComponent],
  providers: [DateFormatPipe],
})
export class KtbEvaluationChartModule {}
