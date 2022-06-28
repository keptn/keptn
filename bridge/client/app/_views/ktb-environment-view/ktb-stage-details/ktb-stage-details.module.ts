import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbStageDetailsComponent } from './ktb-stage-details.component';
import { DtToggleButtonGroupModule } from '@dynatrace/barista-components/toggle-button-group';
import { FlexModule } from '@angular/flex-layout';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbExpandableTileModule } from '../../../_components/ktb-expandable-tile/ktb-expandable-tile.module';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';
import { KtbEvaluationInfoModule } from '../../../_components/ktb-evaluation-info/ktb-evaluation-info.module';
import { KtbLoadingModule } from '../../../_components/ktb-loading/ktb-loading.module';
import { KtbApprovalItemModule } from '../../../_components/ktb-approval-item/ktb-approval-item.module';
import { KtbNoServiceInfoModule } from '../../../_components/ktb-no-service-info/ktb-no-service-info.module';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [KtbStageDetailsComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInfoGroupModule,
    DtToggleButtonGroupModule,
    FlexModule,
    KtbApprovalItemModule,
    KtbExpandableTileModule,
    KtbEvaluationInfoModule,
    KtbLoadingModule,
    KtbNoServiceInfoModule,
    KtbPipeModule,
    RouterModule,
  ],
  exports: [KtbStageDetailsComponent],
})
export class KtbStageDetailsModule {}
