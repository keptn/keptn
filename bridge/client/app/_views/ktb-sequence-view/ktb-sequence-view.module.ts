import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSequenceViewComponent } from './ktb-sequence-view.component';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MomentModule } from 'ngx-moment';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbSelectableTileModule } from '../../_components/ktb-selectable-tile/ktb-selectable-tile.module';
import { KtbSequenceStateInfoModule } from '../../_components/ktb-sequence-state-info/ktb-sequence-state-info.module';
import { KtbRootEventsListComponent } from './ktb-root-events-list/ktb-root-events-list.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtQuickFilterModule } from '@dynatrace/barista-components/quick-filter';
import { KtbLoadingModule } from '../../_components/ktb-loading/ktb-loading.module';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtAlertModule } from '@dynatrace/barista-components/alert';
import { KtbSequenceTimelineComponent } from './ktb-sequence-timeline/ktb-sequence-timeline.component';
import { KtbSequenceTasksListComponent } from './ktb-sequence-tasks-list/ktb-sequence-tasks-list.component';
import { KtbTaskItemComponent, KtbTaskItemDetailDirective } from './ktb-sequence-tasks-list/ktb-task-item.component';
import { MatDialogModule } from '@angular/material/dialog';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { KtbApprovalItemModule } from '../../_components/ktb-approval-item/ktb-approval-item.module';
import { KtbEvaluationInfoModule } from '../../_components/ktb-evaluation-info/ktb-evaluation-info.module';
import { KtbExpandableTileModule } from '../../_components/ktb-expandable-tile/ktb-expandable-tile.module';
import { KtbConfirmationDialogModule } from '../../_components/_dialogs/ktb-confirmation-dialog/ktb-confirmation-dialog.module';
import { KtbSequenceControlsComponent } from './ktb-sequence-controls/ktb-sequence-controls.component';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    component: KtbSequenceViewComponent,
  },
  { path: ':shkeptncontext', component: KtbSequenceViewComponent },
  { path: ':shkeptncontext/event/:eventId', component: KtbSequenceViewComponent },
  { path: ':shkeptncontext/stage/:stage', component: KtbSequenceViewComponent },
];

@NgModule({
  declarations: [
    KtbSequenceViewComponent,
    KtbRootEventsListComponent,
    KtbSequenceControlsComponent,
    KtbSequenceTasksListComponent,
    KtbSequenceTimelineComponent,
    KtbTaskItemComponent,
    KtbTaskItemDetailDirective,
  ],
  imports: [
    CommonModule,
    DtAlertModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInfoGroupModule,
    DtQuickFilterModule,
    DtShowMoreModule,
    DtTagModule,
    FlexLayoutModule,
    KtbApprovalItemModule,
    KtbConfirmationDialogModule,
    KtbEvaluationInfoModule,
    KtbExpandableTileModule,
    KtbLoadingModule,
    KtbPipeModule,
    KtbSelectableTileModule,
    KtbSequenceStateInfoModule,
    MatDialogModule,
    MomentModule,
    RouterModule.forChild(routes),
  ],
})
export class KtbSequenceViewModule {}
