import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbServiceDetailsComponent } from './ktb-service-details.component';
import { RouterModule } from '@angular/router';
import { MatDialogModule } from '@angular/material/dialog';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { KtbDeploymentStageTimelineModule } from '../ktb-deployment-stage-timeline/ktb-deployment-stage-timeline.module';
import { KtbSequenceListModule } from '../ktb-sequence-list/ktb-sequence-list.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbEventItemModule } from '../ktb-event-item/ktb-event-item.module';
import { KtbEvaluationDetailsModule } from '../ktb-evaluation-details/ktb-evaluation-details.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { FlexLayoutModule } from '@angular/flex-layout';

@NgModule({
  declarations: [KtbServiceDetailsComponent],
  imports: [
    CommonModule,
    RouterModule,
    FlexLayoutModule,
    MatDialogModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInfoGroupModule,
    DtTagModule,
    KtbDeploymentStageTimelineModule,
    KtbEvaluationDetailsModule,
    KtbEventItemModule,
    KtbSequenceListModule,
  ],
  exports: [KtbServiceDetailsComponent],
})
export class KtbServiceDetailsModule {}
