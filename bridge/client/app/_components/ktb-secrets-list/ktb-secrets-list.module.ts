import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSecretsListComponent } from './ktb-secrets-list.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbDeleteConfirmationModule } from '../_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [KtbSecretsListComponent],
  imports: [
    CommonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtTableModule,
    KtbPipeModule,
    KtbDeleteConfirmationModule,
    FlexLayoutModule,
    DtButtonModule,
    RouterModule,
  ],
  exports: [KtbSecretsListComponent],
})
export class KtbSecretsListModule {}
