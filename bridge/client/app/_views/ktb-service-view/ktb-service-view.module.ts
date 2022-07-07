import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbServiceViewComponent } from './ktb-service-view.component';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { KtbServiceViewRoutingModule } from './ktb-service-view-routing.module';
import { KtbNoServiceInfoModule } from '../../_components/ktb-no-service-info/ktb-no-service-info.module';
import { KtbExpandableTileModule } from '../../_components/ktb-expandable-tile/ktb-expandable-tile.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbDeploymentListComponent } from './ktb-deployment-list/ktb-deployment-list.component';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { KtbServiceDetailsComponent } from './ktb-service-details/ktb-service-details.component';
import { FlexModule } from '@angular/flex-layout';
import { MatDialogModule } from '@angular/material/dialog';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbEvaluationDetailsModule } from '../../_components/ktb-evaluation-details/ktb-evaluation-details.module';
import { KtbEventItemModule } from '../../_components/ktb-event-item/ktb-event-item.module';
import { KtbLoadingModule } from '../../_components/ktb-loading/ktb-loading.module';
import { KtbSequenceListComponent } from './ktb-sequence-list/ktb-sequence-list.component';
import { MomentModule } from 'ngx-moment';
import { KtbDeploymentStageTimelineComponent } from './ktb-deployment-stage-timeline/ktb-deployment-stage-timeline.component';
import { KtbStageBadgeModule } from '../../_components/ktb-stage-badge/ktb-stage-badge.module';

@NgModule({
  declarations: [
    KtbServiceViewComponent,
    KtbDeploymentListComponent,
    KtbServiceDetailsComponent,
    KtbSequenceListComponent,
    KtbDeploymentStageTimelineComponent,
  ],
  imports: [
    CommonModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInfoGroupModule,
    DtTableModule,
    DtTagModule,
    FlexModule,
    KtbEvaluationDetailsModule,
    KtbEventItemModule,
    KtbExpandableTileModule,
    KtbLoadingModule,
    KtbNoServiceInfoModule,
    KtbPipeModule,
    KtbServiceViewRoutingModule,
    KtbStageBadgeModule,
    MatDialogModule,
    MomentModule,
  ],
})
export class KtbServiceViewModule {}
