import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbKeptnServicesListComponent } from './ktb-keptn-services-list.component';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { KtbExpandableTileModule } from '../ktb-expandable-tile/ktb-expandable-tile.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { KtbUniformSubscriptionsModule } from '../ktb-uniform-subscriptions/ktb-uniform-subscriptions.module';
import { KtbUniformRegistrationLogsModule } from '../ktb-uniform-registration-logs/ktb-uniform-registration-logs.module';

@NgModule({
  declarations: [KtbKeptnServicesListComponent],
  imports: [
    CommonModule,
    DtTableModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtOverlayModule,
    FlexLayoutModule,
    KtbExpandableTileModule,
    KtbLoadingModule,
    KtbPipeModule,
    KtbUniformRegistrationLogsModule,
    KtbUniformSubscriptionsModule,
  ],
  exports: [KtbKeptnServicesListComponent],
})
export class KtbKeptnServicesListModule {}
