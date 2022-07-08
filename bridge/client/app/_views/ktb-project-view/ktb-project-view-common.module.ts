import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbProjectViewComponent } from './ktb-project-view.component';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { RouterModule } from '@angular/router';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { FlexModule } from '@angular/flex-layout';
import { KtbLoadingModule } from '../../_components/ktb-loading/ktb-loading.module';
import { DtMenuModule } from '@dynatrace/barista-components/menu';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { KtbCreateSecretFormModule } from '../../_components/ktb-create-secret-form/ktb-create-secret-form.module';
import { KtbCreateServiceModule } from '../../_components/ktb-create-service/ktb-create-service.module';
import { KtbEditServiceModule } from '../../_components/ktb-edit-service/ktb-edit-service.module';
import { KtbProjectSettingsModule } from '../../_components/ktb-project-settings/ktb-project-settings.module';
import { KtbSecretsListModule } from '../../_components/ktb-secrets-list/ktb-secrets-list.module';
import { KtbServiceSettingsModule } from '../../_components/ktb-service-settings/ktb-service-settings.module';

@NgModule({
  declarations: [KtbProjectViewComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtEmptyStateModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtMenuModule,
    DtOverlayModule,
    FlexModule,
    KtbLoadingModule,
    KtbPipeModule,
    RouterModule,
    KtbCreateSecretFormModule,
    KtbCreateServiceModule,
    KtbEditServiceModule,
    KtbProjectSettingsModule,
    KtbSecretsListModule,
    KtbServiceSettingsModule,
  ],
})
export class KtbProjectViewCommonModule {}
