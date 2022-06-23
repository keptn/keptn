import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbPayloadViewerComponent } from './ktb-payload-viewer.component';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { DtAlertModule } from '@dynatrace/barista-components/alert';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';

@NgModule({
  declarations: [KtbPayloadViewerComponent],
  imports: [
    CommonModule,
    DtAlertModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbLoadingModule,
    MatDialogModule,
  ],
  exports: [KtbPayloadViewerComponent],
})
export class KtbPayloadViewerModule {}
