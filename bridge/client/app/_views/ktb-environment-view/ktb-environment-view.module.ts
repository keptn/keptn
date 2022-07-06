import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { DtToggleButtonGroupModule } from '@dynatrace/barista-components/toggle-button-group';
import { FlexModule } from '@angular/flex-layout';
import { KtbApprovalItemModule } from '../../_components/ktb-approval-item/ktb-approval-item.module';
import { KtbExpandableTileModule } from '../../_components/ktb-expandable-tile/ktb-expandable-tile.module';
import { KtbEvaluationInfoModule } from '../../_components/ktb-evaluation-info/ktb-evaluation-info.module';
import { KtbLoadingModule } from '../../_components/ktb-loading/ktb-loading.module';
import { KtbNoServiceInfoModule } from '../../_components/ktb-no-service-info/ktb-no-service-info.module';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbStageDetailsComponent } from './ktb-stage-details/ktb-stage-details.component';
import { KtbStageOverviewComponent } from './ktb-stage-overview/ktb-stage-overview.component';
import { KtbServicesListComponent } from './ktb-stage-overview/ktb-services-list/ktb-services-list.component';
import { KtbTriggerSequenceComponent } from './ktb-stage-overview/ktb-trigger-sequence/ktb-trigger-sequence.component';
import { DtFilterFieldModule } from '@dynatrace/barista-components/filter-field';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtRadioModule } from '@dynatrace/barista-components/radio';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { FormsModule } from '@angular/forms';
import { KtbDateInputModule } from '../../_components/ktb-date-input/ktb-date-input.module';
import { KtbSelectableTileModule } from '../../_components/ktb-selectable-tile/ktb-selectable-tile.module';
import { KtbEnvironmentViewRoutingModule } from './ktb-environment-view-routing.module';
import { KtbEnvironmentViewComponent } from './ktb-environment-view.component';

@NgModule({
  declarations: [
    KtbEnvironmentViewComponent,
    KtbStageDetailsComponent,
    KtbStageOverviewComponent,
    KtbServicesListComponent,
    KtbTriggerSequenceComponent,
  ],
  imports: [
    CommonModule,
    DtButtonModule,
    DtFilterFieldModule,
    DtFormFieldModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInfoGroupModule,
    DtInputModule,
    DtOverlayModule,
    DtRadioModule,
    DtSelectModule,
    DtShowMoreModule,
    DtTableModule,
    DtToggleButtonGroupModule,
    FlexModule,
    FormsModule,
    KtbApprovalItemModule,
    KtbDateInputModule,
    KtbExpandableTileModule,
    KtbEvaluationInfoModule,
    KtbLoadingModule,
    KtbNoServiceInfoModule,
    KtbPipeModule,
    KtbSelectableTileModule,
    KtbEnvironmentViewRoutingModule,
  ],
})
export class KtbEnvironmentViewModule {}
