import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MomentModule } from 'ngx-moment';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbSelectableTileModule } from '../ktb-selectable-tile/ktb-selectable-tile.module';
import { KtbSequenceControlsModule } from '../ktb-sequence-controls/ktb-sequence-controls.module';
import { KtbSequenceStateInfoModule } from '../ktb-sequence-state-info/ktb-sequence-state-info.module';
import { KtbRootEventsListComponent } from './ktb-root-events-list.component';

@NgModule({
  declarations: [KtbRootEventsListComponent],
  imports: [
    CommonModule,
    FlexLayoutModule,
    MomentModule,
    KtbPipeModule,
    KtbSelectableTileModule,
    KtbSequenceControlsModule,
    KtbSequenceStateInfoModule,
  ],
  exports: [KtbRootEventsListComponent],
})
export class KtbRootEventsListModule {}
