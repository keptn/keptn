import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FlexModule } from '@angular/flex-layout';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtConfirmationDialogModule } from '@dynatrace/barista-components/confirmation-dialog';
import { DtContextDialogModule } from '@dynatrace/barista-components/context-dialog';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { DtSwitchModule } from '@dynatrace/barista-components/switch';
import { DtTopBarNavigationModule } from '@dynatrace/barista-components/top-bar-navigation';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbCopyToClipboardModule } from '../ktb-copy-to-clipboard/ktb-copy-to-clipboard.module';
import { KtbAppHeaderComponent } from './ktb-app-header.component';
import { KtbUserComponent } from './ktb-user/ktb-user.component';

@NgModule({
  declarations: [KtbAppHeaderComponent, KtbUserComponent],
  imports: [
    CommonModule,
    HttpClientModule,
    DtButtonModule,
    DtConfirmationDialogModule,
    DtContextDialogModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtSelectModule,
    DtSwitchModule,
    DtTopBarNavigationModule,
    FlexModule,
    FormsModule,
    KtbCopyToClipboardModule,
    KtbPipeModule,
    RouterModule,
  ],
  exports: [KtbAppHeaderComponent],
})
export class KtbAppHeaderModule {}
