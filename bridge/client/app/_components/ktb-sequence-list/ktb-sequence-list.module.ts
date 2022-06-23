import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSequenceListComponent } from './ktb-sequence-list.component';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { MomentModule } from 'ngx-moment';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { RouterModule } from '@angular/router';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';

@NgModule({
  declarations: [KtbSequenceListComponent],
  imports: [
    CommonModule,
    RouterModule,
    FlexLayoutModule,
    MomentModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtTableModule,
    KtbLoadingModule,
    KtbPipeModule,
  ],
  exports: [KtbSequenceListComponent],
})
export class KtbSequenceListModule {}
