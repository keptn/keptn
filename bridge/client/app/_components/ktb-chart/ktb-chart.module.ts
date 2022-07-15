import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbChartComponent } from './ktb-chart.component';
import { KtbTestChartComponent } from './testing/ktb-test-chart.component';
import { FlexModule } from '@angular/flex-layout';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';

@NgModule({
  declarations: [KtbChartComponent, KtbTestChartComponent], // add KtbTestChartComponent for testing
  imports: [
    CommonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtButtonModule,
    FlexModule,
    DtKeyValueListModule,
    KtbPipeModule,
  ],
  exports: [KtbChartComponent, KtbTestChartComponent], // add KtbTestChartComponent for testing,
})
export class KtbChartModule {}
