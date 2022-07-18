import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtAlertModule } from '@dynatrace/barista-components/alert';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbNotificationComponent } from './ktb-notification.component';
import { KtbNotificationBarComponent } from './ktb-notification-bar.component';

@NgModule({
  declarations: [KtbNotificationComponent, KtbNotificationBarComponent],
  imports: [
    CommonModule,
    FlexLayoutModule,
    DtAlertModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtButtonModule,
    KtbPipeModule,
  ],
  exports: [KtbNotificationBarComponent],
})
export class KtbNotificationModule {}
