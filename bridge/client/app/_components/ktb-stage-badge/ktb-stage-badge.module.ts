import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbStageBadgeComponent } from './ktb-stage-badge.component';
import { FlexModule } from '@angular/flex-layout';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { KtbEvaluationInfoModule } from '../ktb-evaluation-info/ktb-evaluation-info.module';

@NgModule({
  declarations: [KtbStageBadgeComponent],
  imports: [CommonModule, FlexModule, DtTagModule, KtbEvaluationInfoModule],
  exports: [KtbStageBadgeComponent],
})
export class KtbStageBadgeModule {}
