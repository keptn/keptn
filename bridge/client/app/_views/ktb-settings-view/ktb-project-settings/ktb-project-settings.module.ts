import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtCheckboxModule } from '@dynatrace/barista-components/checkbox';
import { DtConfirmationDialogModule } from '@dynatrace/barista-components/confirmation-dialog';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtRadioModule } from '@dynatrace/barista-components/radio';
import { DtSwitchModule } from '@dynatrace/barista-components/switch';
import { KtbDragAndDropModule } from '../../../_directives/ktb-drag-and-drop/ktb-drag-and-drop.module';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';
import { KtbErrorViewModule } from '../../ktb-error-view/ktb-error-view.module';
import { KtbCertificateInputModule } from '../../../_components/ktb-certificate-input/ktb-certificate-input.module';
import { KtbDangerZoneModule } from '../../../_components/ktb-danger-zone/ktb-danger-zone.module';
import { KtbLoadingModule } from '../../../_components/ktb-loading/ktb-loading.module';
import { KtbProxyInputModule } from '../../../_components/ktb-proxy-input/ktb-proxy-input.module';
import { KtbSshKeyInputModule } from '../../../_components/ktb-ssh-key-input/ktb-ssh-key-input.module';
import { KtbProjectCreateMessageComponent } from './ktb-project-create-message/ktb-project-create-message.component';
import { KtbProjectSettingsGitExtendedComponent } from './ktb-project-settings-git-extended/ktb-project-settings-git-extended.component';
import { KtbProjectSettingsGitHttpsComponent } from './ktb-project-settings-git-https/ktb-project-settings-git-https.component';
import { KtbProjectSettingsGitSshInputComponent } from './ktb-project-settings-git-ssh-input/ktb-project-settings-git-ssh-input.component';
import { KtbProjectSettingsGitSshComponent } from './ktb-project-settings-git-ssh/ktb-project-settings-git-ssh.component';
import { KtbProjectSettingsGitComponent } from './ktb-project-settings-git/ktb-project-settings-git.component';
import { KtbProjectSettingsShipyardComponent } from './ktb-project-settings-shipyard/ktb-project-settings-shipyard.component';
import { KtbProjectSettingsComponent } from './ktb-project-settings.component';

@NgModule({
  declarations: [
    KtbProjectSettingsComponent,
    KtbProjectCreateMessageComponent,
    KtbProjectSettingsGitComponent,
    KtbProjectSettingsGitExtendedComponent,
    KtbProjectSettingsGitHttpsComponent,
    KtbProjectSettingsGitSshComponent,
    KtbProjectSettingsGitSshInputComponent,
    KtbProjectSettingsShipyardComponent,
  ],
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    FlexLayoutModule,
    DtButtonModule,
    DtCheckboxModule,
    DtConfirmationDialogModule,
    DtFormFieldModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInputModule,
    DtOverlayModule,
    DtRadioModule,
    DtSwitchModule,
    KtbCertificateInputModule,
    KtbDangerZoneModule,
    KtbDragAndDropModule,
    KtbErrorViewModule,
    KtbLoadingModule,
    KtbPipeModule,
    KtbProxyInputModule,
    KtbSshKeyInputModule,
  ],
  exports: [KtbProjectSettingsComponent],
})
export class KtbProjectSettingsModule {}
