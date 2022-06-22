import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbServiceSettingsComponent } from './ktb-service-settings.component';
import { RouterModule } from '@angular/router';
import { KtbServiceSettingsListComponent } from './ktb-service-settings-list/ktb-service-settings-list.component';
import { KtbServiceSettingsOverviewComponent } from './ktb-service-settings-overview/ktb-service-settings-overview.component';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { FlexLayoutModule } from '@angular/flex-layout';

@NgModule({
  declarations: [KtbServiceSettingsComponent, KtbServiceSettingsListComponent, KtbServiceSettingsOverviewComponent],
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
  exports: [KtbServiceSettingsComponent],
})
export class KtbServiceSettingsModule {}
