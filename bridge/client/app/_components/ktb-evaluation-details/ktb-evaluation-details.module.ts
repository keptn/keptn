import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatDialogModule } from '@angular/material/dialog';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtButtonGroupModule } from '@dynatrace/barista-components/button-group';
import { DtChartModule } from '@dynatrace/barista-components/chart';
import { DtConsumptionModule } from '@dynatrace/barista-components/consumption';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { MomentModule } from 'ngx-moment';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbHeatmapModule } from '../ktb-heatmap/ktb-heatmap.module';
import { KtbSliBreakdownModule } from '../ktb-sli-breakdown/ktb-sli-breakdown.module';
import { KtbEvaluationDetailsComponent } from './ktb-evaluation-details.component';

@NgModule({
  declarations: [KtbEvaluationDetailsComponent],
  imports: [
    CommonModule,
    HttpClientModule,
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
  ],
  exports: [KtbEvaluationDetailsComponent],
})
export class KtbEvaluationDetailsModule {}
