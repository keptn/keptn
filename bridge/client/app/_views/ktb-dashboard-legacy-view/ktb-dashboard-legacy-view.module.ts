import { NgModule } from '@angular/core';
import { KtbDashboardLegacyViewComponent } from './ktb-dashboard-legacy-view.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbLoadingModule } from '../../_components/ktb-loading/ktb-loading.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { KtbDashboardLegacyViewRoutingModule } from './ktb-dashboard-legacy-view-routing.module';
import { KtbProjectListComponent } from './ktb-project-list/ktb-project-list.component';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { DtTileModule } from '@dynatrace/barista-components/tile';
import { KtbSequenceStateListModule } from '../../_components/ktb-sequence-state-list/ktb-sequence-state-list.module';
import { KtbProjectTileComponent } from './ktb-project-list/ktb-project-tile.component';

@NgModule({
  declarations: [KtbDashboardLegacyViewComponent, KtbProjectTileComponent, KtbProjectListComponent],
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
    FlexLayoutModule,
    KtbDashboardLegacyViewRoutingModule,
    RouterModule,
    DtTagModule,
    DtTileModule,
    KtbSequenceStateListModule,
  ],
})
export class KtbDashboardLegacyViewModule {}
