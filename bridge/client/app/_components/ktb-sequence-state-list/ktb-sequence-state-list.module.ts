import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { MomentModule } from 'ngx-moment';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbSequenceStateInfoModule } from '../ktb-sequence-state-info/ktb-sequence-state-info.module';
import { KtbSequenceStateListComponent } from './ktb-sequence-state-list.component';

@NgModule({
  declarations: [KtbSequenceStateListComponent],
  imports: [
    CommonModule,
    DtTableModule,
    KtbPipeModule,
    KtbSequenceStateInfoModule,
    MomentModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    FlexLayoutModule,
  ],
  exports: [KtbSequenceStateListComponent],
})
export class KtbSequenceStateListModule {}
