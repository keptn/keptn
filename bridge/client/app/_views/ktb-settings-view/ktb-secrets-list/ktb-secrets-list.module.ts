import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSecretsListComponent } from './ktb-secrets-list.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';
import { KtbDeleteConfirmationModule } from '../../../_components/_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [KtbSecretsListComponent],
  imports: [
    CommonModule,
    FlexLayoutModule,
    RouterModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtTableModule,
    DtButtonModule,
    KtbPipeModule,
    KtbDeleteConfirmationModule,
  ],
  exports: [KtbSecretsListComponent],
})
export class KtbSecretsListModule {}
