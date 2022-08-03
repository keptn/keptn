import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEvaluationDetailsComponent } from './ktb-evaluation-details.component';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonGroupModule } from '@dynatrace/barista-components/button-group';
import { KtbHeatmapModule } from '../ktb-heatmap/ktb-heatmap.module';
import { DtChartModule } from '@dynatrace/barista-components/chart';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { DtConsumptionModule } from '@dynatrace/barista-components/consumption';
import { DateFormatPipe, MomentModule } from 'ngx-moment';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { MatDialogModule } from '@angular/material/dialog';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbSliBreakdownModule } from '../ktb-sli-breakdown/ktb-sli-breakdown.module';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { KtbChartModule } from '../ktb-chart/ktb-chart.module';

@NgModule({
  declarations: [KtbEvaluationDetailsComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtButtonGroupModule,
    DtChartModule,
    DtConsumptionModule,
    DtFormFieldModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInputModule,
    DtKeyValueListModule,
    DtOverlayModule,
    FlexLayoutModule,
    KtbHeatmapModule,
    KtbPipeModule,
    KtbSliBreakdownModule,
    MatDialogModule,
    MomentModule,
    KtbChartModule,
  ],
  exports: [KtbEvaluationDetailsComponent],
  providers: [DateFormatPipe],
})
export class KtbEvaluationDetailsModule {}
