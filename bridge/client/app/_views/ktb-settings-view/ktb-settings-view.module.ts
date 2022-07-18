import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtMenuModule } from '@dynatrace/barista-components/menu';
import { KtbCreateSecretFormModule } from './ktb-create-secret-form/ktb-create-secret-form.module';
import { KtbCreateServiceModule } from './ktb-create-service/ktb-create-service.module';
import { KtbEditServiceModule } from './ktb-edit-service/ktb-edit-service.module';
import { KtbProjectSettingsModule } from './ktb-project-settings/ktb-project-settings.module';
import { KtbSecretsListModule } from './ktb-secrets-list/ktb-secrets-list.module';
import { KtbServiceSettingsModule } from './ktb-service-settings/ktb-service-settings.module';
import { KtbSettingsViewRoutingModule } from './ktb-settings-view-routing.module';
import { KtbSettingsViewComponent } from './ktb-settings-view.component';

@NgModule({
  declarations: [KtbSettingsViewComponent],
  imports: [
    CommonModule,
    FlexLayoutModule,
    RouterModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtMenuModule,
    KtbSettingsViewRoutingModule,
    KtbCreateSecretFormModule,
    KtbCreateServiceModule,
    KtbEditServiceModule,
    KtbProjectSettingsModule,
    KtbSecretsListModule,
    KtbServiceSettingsModule,
  ],
})
export class KtbSettingsViewModule {}
