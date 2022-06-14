import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEventItemComponent, KtbEventItemDetailDirective } from './ktb-event-item.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MomentModule } from 'ngx-moment';
import { RouterModule } from '@angular/router';
import { KtbSelectableTileModule } from '../ktb-selectable-tile/ktb-selectable-tile.module';
import { KtbApprovalItemModule } from '../ktb-approval-item/ktb-approval-item.module';
import { MatDialogModule } from '@angular/material/dialog';

@NgModule({
  declarations: [KtbEventItemComponent, KtbEventItemDetailDirective],
  imports: [
    CommonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtButtonModule,
    FlexLayoutModule,
    KtbApprovalItemModule,
    KtbSelectableTileModule,
    MatDialogModule,
    MomentModule,
    RouterModule,
  ],
  exports: [KtbEventItemComponent, KtbEventItemDetailDirective],
})
export class KtbEventItemModule {}
