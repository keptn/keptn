import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbNotificationComponent } from './ktb-notification.component';
import { KtbNotificationBarComponent } from './ktb-notification-bar.component';
import { DtAlertModule } from '@dynatrace/barista-components/alert';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';

@NgModule({
  declarations: [KtbNotificationComponent, KtbNotificationBarComponent],
  imports: [
    CommonModule,
    DtAlertModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbPipeModule,
  ],
  exports: [KtbNotificationBarComponent],
})
export class KtbNotificationModule {}
