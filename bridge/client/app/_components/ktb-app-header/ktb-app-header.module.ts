import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbUserComponent } from './ktb-user/ktb-user.component';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { FormsModule } from '@angular/forms';
import { KtbAppHeaderComponent } from './ktb-app-header.component';
import { RouterModule } from '@angular/router';
import { DtSwitchModule } from '@dynatrace/barista-components/switch';
import { FlexModule } from '@angular/flex-layout';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { DtConfirmationDialogModule } from '@dynatrace/barista-components/confirmation-dialog';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { DtTopBarNavigationModule } from '@dynatrace/barista-components/top-bar-navigation';
import { DtContextDialogModule } from '@dynatrace/barista-components/context-dialog';
import { KtbCopyToClipboardModule } from '../ktb-copy-to-clipboard/ktb-copy-to-clipboard.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { HttpClientModule } from '@angular/common/http';

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
