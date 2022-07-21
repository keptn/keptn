import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FlexModule } from '@angular/flex-layout';
import { MatDialogModule } from '@angular/material/dialog';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtMenuModule } from '@dynatrace/barista-components/menu';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbProjectViewComponent } from './ktb-project-view.component';
import { KtbLoadingModule } from '../../_components/ktb-loading/ktb-loading.module';

@NgModule({
  declarations: [KtbProjectViewComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtEmptyStateModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtMenuModule,
    DtOverlayModule,
    FlexModule,
    MatDialogModule,
    KtbLoadingModule,
    KtbPipeModule,
    RouterModule,
  ],
})
export class KtbProjectViewCommonModule {}
