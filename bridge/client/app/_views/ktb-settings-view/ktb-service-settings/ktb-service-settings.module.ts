import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbLoadingModule } from '../../../_components/ktb-loading/ktb-loading.module';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';
import { KtbServiceSettingsListComponent } from './ktb-service-settings-list.component';
import { KtbServiceSettingsOverviewComponent } from './ktb-service-settings-overview.component';

@NgModule({
  declarations: [KtbServiceSettingsListComponent, KtbServiceSettingsOverviewComponent],
  imports: [
    CommonModule,
    RouterModule,
    FlexLayoutModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtTableModule,
    KtbLoadingModule,
    KtbPipeModule,
  ],
  exports: [KtbServiceSettingsListComponent, KtbServiceSettingsOverviewComponent],
})
export class KtbServiceSettingsModule {}
