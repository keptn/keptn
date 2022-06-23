import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbStageOverviewComponent } from './ktb-stage-overview.component';
import { DtFilterFieldModule } from '@dynatrace/barista-components/filter-field';
import { FlexModule } from '@angular/flex-layout';
import { KtbTriggerSequenceModule } from '../../../_components/ktb-trigger-sequence/ktb-trigger-sequence.module';
import { KtbSelectableTileModule } from '../../../_components/ktb-selectable-tile/ktb-selectable-tile.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { RouterModule } from '@angular/router';
import { KtbServicesListComponent } from './ktb-services-list/ktb-services-list.component';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';
import { KtbLoadingModule } from '../../../_components/ktb-loading/ktb-loading.module';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { KtbEvaluationInfoModule } from '../../../_components/ktb-evaluation-info/ktb-evaluation-info.module';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { KtbNoServiceInfoModule } from '../../../_components/ktb-no-service-info/ktb-no-service-info.module';

@NgModule({
  declarations: [KtbStageOverviewComponent, KtbServicesListComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtFilterFieldModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtOverlayModule,
    DtShowMoreModule,
    DtTableModule,
    FlexModule,
    KtbEvaluationInfoModule,
    KtbLoadingModule,
    KtbNoServiceInfoModule,
    KtbPipeModule,
    KtbSelectableTileModule,
    KtbTriggerSequenceModule,
    RouterModule,
  ],
  exports: [KtbStageOverviewComponent],
})
export class KtbStageOverviewModule {}
