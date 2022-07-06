import { NgModule } from '@angular/core';
import { KtbDashboardLegacyViewComponent } from './ktb-dashboard-legacy-view.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbLoadingModule } from '../../_components/ktb-loading/ktb-loading.module';
import { FlexModule } from '@angular/flex-layout';
import { CommonModule } from '@angular/common';
import { KtbProjectListModule } from '../../_components/ktb-project-list/ktb-project-list.module';
import { RouterModule } from '@angular/router';
import { KtbDashboardLegacyViewRoutingModule } from './ktb-dashboard-legacy-view-routing.module';

@NgModule({
  declarations: [KtbDashboardLegacyViewComponent],
  imports: [
    CommonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtEmptyStateModule,
    DtInfoGroupModule,
    KtbPipeModule,
    DtButtonModule,
    KtbLoadingModule,
    FlexModule,
    KtbProjectListModule,
    KtbDashboardLegacyViewRoutingModule,
    RouterModule,
  ],
})
export class KtbDashboardLegacyViewModule {}
