import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEvaluationInfoComponent } from './ktb-evaluation-info.component';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { FlexModule } from '@angular/flex-layout';
import { KtbEvaluationDetailsModule } from '../ktb-evaluation-details/ktb-evaluation-details.module';
import { KtbEvaluationBadgeModule } from '../ktb-evaluation-badge/ktb-evaluation-badge.module';

@NgModule({
  declarations: [KtbEvaluationInfoComponent],
  imports: [CommonModule, DtTagModule, FlexModule, KtbEvaluationDetailsModule, KtbEvaluationBadgeModule],
  exports: [KtbEvaluationInfoComponent],
})
export class KtbEvaluationInfoModule {}
