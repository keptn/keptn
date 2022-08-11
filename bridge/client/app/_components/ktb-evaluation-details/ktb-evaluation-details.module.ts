import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEvaluationDetailsComponent } from './ktb-evaluation-details.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtConsumptionModule } from '@dynatrace/barista-components/consumption';
import { DateFormatPipe, MomentModule } from 'ngx-moment';
import { KtbSliBreakdownModule } from '../ktb-sli-breakdown/ktb-sli-breakdown.module';
import { KtbEvaluationChartModule } from './ktb-evaluation-chart/ktb-evaluation-chart.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { FlexModule } from '@angular/flex-layout';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { MatDialogModule } from '@angular/material/dialog';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';

@NgModule({
  declarations: [KtbEvaluationDetailsComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtConsumptionModule,
    DtFormFieldModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInputModule,
    DtOverlayModule,
    FlexModule,
    KtbEvaluationChartModule,
    KtbPipeModule,
    KtbSliBreakdownModule,
    MatDialogModule,
    MomentModule,
  ],
  exports: [KtbEvaluationDetailsComponent],
  providers: [DateFormatPipe],
})
export class KtbEvaluationDetailsModule {}
