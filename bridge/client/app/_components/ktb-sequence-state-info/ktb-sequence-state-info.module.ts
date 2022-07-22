import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { RouterModule } from '@angular/router';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { KtbStageBadgeModule } from '../ktb-stage-badge/ktb-stage-badge.module';
import { KtbSequenceStateInfoComponent } from './ktb-sequence-state-info.component';

@NgModule({
  declarations: [KtbSequenceStateInfoComponent],
  imports: [
    CommonModule,
    RouterModule,
    HttpClientModule,
    FlexLayoutModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbLoadingModule,
    KtbStageBadgeModule,
  ],
  exports: [KtbSequenceStateInfoComponent],
})
export class KtbSequenceStateInfoModule {}
