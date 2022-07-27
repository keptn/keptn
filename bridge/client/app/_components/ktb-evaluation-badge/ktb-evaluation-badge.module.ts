import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEvaluationBadgeComponent } from './ktb-evaluation-badge.component';
import { KtbEvaluationDetailsModule } from '../ktb-evaluation-details/ktb-evaluation-details.module';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';

@NgModule({
  declarations: [KtbEvaluationBadgeComponent],
  imports: [CommonModule, KtbEvaluationDetailsModule, KtbLoadingModule, KtbPipeModule],
  exports: [KtbEvaluationBadgeComponent],
})
export class KtbEvaluationBadgeModule {}
