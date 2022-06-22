import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbServicesListComponent } from './ktb-services-list.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { RouterModule } from '@angular/router';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { KtbEvaluationInfoModule } from '../ktb-evaluation-info/ktb-evaluation-info.module';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { KtbNoServiceInfoModule } from '../ktb-no-service-info/ktb-no-service-info.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';

@NgModule({
  declarations: [KtbServicesListComponent],
  imports: [
    CommonModule,
    RouterModule,
    FlexLayoutModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtOverlayModule,
    DtShowMoreModule,
    DtTableModule,
    KtbEvaluationInfoModule,
    KtbLoadingModule,
    KtbNoServiceInfoModule,
    KtbPipeModule,
  ],
  exports: [KtbServicesListComponent],
})
export class KtbServicesListModule {}
