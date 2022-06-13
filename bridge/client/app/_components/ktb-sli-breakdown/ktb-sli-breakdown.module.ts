import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSliBreakdownComponent } from './ktb-sli-breakdown.component';
import { KtbSliBreakdownCriteriaItemComponent } from './ktb-sli-breakdown-criteria-item.component';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { FlexModule } from '@angular/flex-layout';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtButtonModule } from '@dynatrace/barista-components/button';

@NgModule({
  declarations: [KtbSliBreakdownComponent, KtbSliBreakdownCriteriaItemComponent],
  imports: [
    CommonModule,
    DtTableModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbPipeModule,
    FlexModule,
    DtOverlayModule,
    DtButtonModule,
  ],
  exports: [KtbSliBreakdownComponent],
})
export class KtbSliBreakdownModule {}
