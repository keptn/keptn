import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { RouterModule } from '@angular/router';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { DtTileModule } from '@dynatrace/barista-components/tile';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbSequenceStateListModule } from '../ktb-sequence-state-list/ktb-sequence-state-list.module';
import { KtbProjectListComponent } from './ktb-project-list.component';
import { KtbProjectTileComponent } from './ktb-project-tile.component';

@NgModule({
  declarations: [KtbProjectListComponent, KtbProjectTileComponent],
  imports: [
    CommonModule,
    RouterModule,
    FlexLayoutModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtTagModule,
    DtTileModule,
    KtbPipeModule,
    KtbSequenceStateListModule,
  ],
  exports: [KtbProjectListComponent, KtbProjectTileComponent],
})
export class KtbProjectListModule {}
