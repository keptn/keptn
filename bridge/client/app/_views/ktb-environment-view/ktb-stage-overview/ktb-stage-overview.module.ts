import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbStageOverviewComponent } from './ktb-stage-overview.component';
import { DtFilterFieldModule } from '@dynatrace/barista-components/filter-field';
import { FlexModule } from '@angular/flex-layout';
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
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtRadioModule } from '@dynatrace/barista-components/radio';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { FormsModule } from '@angular/forms';
import { KtbDateInputModule } from '../../../_components/ktb-date-input/ktb-date-input.module';
import { KtbTriggerSequenceComponent } from './ktb-trigger-sequence/ktb-trigger-sequence.component';

@NgModule({
  declarations: [KtbStageOverviewComponent, KtbServicesListComponent, KtbTriggerSequenceComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtFilterFieldModule,
    DtFormFieldModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInputModule,
    DtOverlayModule,
    DtRadioModule,
    DtSelectModule,
    DtShowMoreModule,
    DtTableModule,
    FlexModule,
    FormsModule,
    KtbDateInputModule,
    KtbEvaluationInfoModule,
    KtbLoadingModule,
    KtbNoServiceInfoModule,
    KtbPipeModule,
    KtbSelectableTileModule,
    RouterModule,
  ],
  exports: [KtbStageOverviewComponent],
})
export class KtbStageOverviewModule {}
