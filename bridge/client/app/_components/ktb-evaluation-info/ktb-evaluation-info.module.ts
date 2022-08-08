import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEvaluationInfoComponent } from './ktb-evaluation-info.component';
import { FlexModule } from '@angular/flex-layout';
import { KtbEvaluationBadgeModule } from '../ktb-evaluation-badge/ktb-evaluation-badge.module';

@NgModule({
  declarations: [KtbEvaluationInfoComponent],
  imports: [CommonModule, FlexModule, KtbEvaluationBadgeModule],
  exports: [KtbEvaluationInfoComponent],
})
export class KtbEvaluationInfoModule {}
