import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbNotificationComponent } from './ktb-notification.component';
import { KtbNotificationBarComponent } from './ktb-notification-bar.component';
import { DtAlertModule } from '@dynatrace/barista-components/alert';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';

@NgModule({
  declarations: [KtbNotificationComponent, KtbNotificationBarComponent],
  imports: [
    CommonModule,
    BrowserAnimationsModule,
    DtAlertModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbPipeModule,
    FlexLayoutModule,
    DtButtonModule,
  ],
  exports: [KtbNotificationBarComponent],
})
export class KtbNotificationModule {}
