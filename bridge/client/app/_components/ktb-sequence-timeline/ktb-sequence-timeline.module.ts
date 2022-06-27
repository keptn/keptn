import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSequenceTimelineComponent } from './ktb-sequence-timeline.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';

@NgModule({
  declarations: [KtbSequenceTimelineComponent],
  imports: [
    CommonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbLoadingModule,
    FlexLayoutModule,
    DtButtonModule,
  ],
  exports: [KtbSequenceTimelineComponent],
})
export class KtbSequenceTimelineModule {}
