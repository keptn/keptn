import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbApprovalItemComponent } from './ktb-approval-item.component';
import { FlexModule } from '@angular/flex-layout';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { KtbEvaluationBadgeModule } from '../ktb-evaluation-badge/ktb-evaluation-badge.module';

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
    KtbEvaluationBadgeModule,
    KtbLoadingModule,
  ],
  exports: [KtbApprovalItemComponent],
})
export class KtbApprovalItemModule {}
