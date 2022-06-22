import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSequenceTasksListComponent } from './ktb-sequence-tasks-list.component';
import { KtbTaskItemComponent, KtbTaskItemDetailDirective } from './ktb-task-item.component';
import { KtbExpandableTileModule } from '../ktb-expandable-tile/ktb-expandable-tile.module';
import { MomentModule } from 'ngx-moment';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbEvaluationInfoModule } from '../ktb-evaluation-info/ktb-evaluation-info.module';
import { MatDialogModule } from '@angular/material/dialog';
import { RouterModule } from '@angular/router';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { KtbApprovalItemModule } from '../ktb-approval-item/ktb-approval-item.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { FlexLayoutModule } from '@angular/flex-layout';

@NgModule({
  declarations: [KtbSequenceTasksListComponent, KtbTaskItemComponent, KtbTaskItemDetailDirective],
  imports: [
    CommonModule,
    RouterModule,
    FlexLayoutModule,
    MomentModule,
    MatDialogModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtTagModule,
    KtbApprovalItemModule,
    KtbEvaluationInfoModule,
    KtbExpandableTileModule,
    KtbLoadingModule,
  ],
  exports: [KtbSequenceTasksListComponent],
})
export class KtbSequenceTasksListModule {}
