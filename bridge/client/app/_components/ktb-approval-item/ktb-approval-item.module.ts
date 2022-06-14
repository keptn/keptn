import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbApprovalItemComponent } from './ktb-approval-item.component';
import { FlexModule } from '@angular/flex-layout';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbEvaluationDetailsModule } from '../ktb-evaluation-details/ktb-evaluation-details.module';

@NgModule({
  declarations: [KtbApprovalItemComponent],
  imports: [
    CommonModule,
    FlexModule,
    DtTagModule,
    DtOverlayModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtButtonModule,
    KtbEvaluationDetailsModule,
  ],
  exports: [KtbApprovalItemComponent],
})
export class KtbApprovalItemModule {}
