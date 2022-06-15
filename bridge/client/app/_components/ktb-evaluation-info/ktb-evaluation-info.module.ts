import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEvaluationInfoComponent } from './ktb-evaluation-info.component';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { FlexModule } from '@angular/flex-layout';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbEvaluationDetailsModule } from '../ktb-evaluation-details/ktb-evaluation-details.module';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';

@NgModule({
  declarations: [KtbEvaluationInfoComponent],
  imports: [CommonModule, DtTagModule, FlexModule, KtbPipeModule, KtbEvaluationDetailsModule, KtbLoadingModule],
  exports: [KtbEvaluationInfoComponent],
})
export class KtbEvaluationInfoModule {}
