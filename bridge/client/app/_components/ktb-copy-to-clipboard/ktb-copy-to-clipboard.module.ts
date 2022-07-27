import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbCopyToClipboardComponent } from './ktb-copy-to-clipboard.component';
import { DtCopyToClipboardModule } from '@dynatrace/barista-components/copy-to-clipboard';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { FlexModule } from '@angular/flex-layout';

@NgModule({
  declarations: [KtbCopyToClipboardComponent],
  imports: [
    CommonModule,
    DtCopyToClipboardModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtButtonModule,
    DtInputModule,
    FlexModule,
  ],
  exports: [KtbCopyToClipboardComponent],
})
export class KtbCopyToClipboardModule {}
