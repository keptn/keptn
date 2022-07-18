import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbChartComponent } from './ktb-chart.component';
import { FlexModule } from '@angular/flex-layout';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';

@NgModule({
  declarations: [KtbChartComponent], // add KtbTestChartComponent for testing
  imports: [
    CommonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    FlexModule,
    DtKeyValueListModule,
    KtbPipeModule,
  ],
  exports: [KtbChartComponent], // add KtbTestChartComponent for testing,
})
export class KtbChartModule {}
