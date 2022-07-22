import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FlexModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtCopyToClipboardModule } from '@dynatrace/barista-components/copy-to-clipboard';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { KtbCopyToClipboardComponent } from './ktb-copy-to-clipboard.component';

@NgModule({
  declarations: [KtbCopyToClipboardComponent],
  imports: [
    CommonModule,
    HttpClientModule,
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
