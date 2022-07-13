import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { KtbSettingsViewComponent } from './ktb-settings-view.component';
import { KtbProjectSettingsComponent } from './ktb-project-settings/ktb-project-settings.component';
import { PendingChangesGuard } from '../../_guards/pending-changes.guard';
import { KtbCreateServiceComponent } from './ktb-create-service/ktb-create-service.component';
import { KtbEditServiceComponent } from './ktb-edit-service/ktb-edit-service.component';
import { KtbServiceSettingsOverviewComponent } from './ktb-service-settings/ktb-service-settings-overview.component';
import { KtbSecretsListComponent } from './ktb-secrets-list/ktb-secrets-list.component';
import { KtbCreateSecretFormComponent } from './ktb-create-secret-form/ktb-create-secret-form.component';

const routes: Routes = [
  {
    path: '',
    component: KtbSettingsViewComponent,
    children: [
      { path: 'project', component: KtbProjectSettingsComponent, canDeactivate: [PendingChangesGuard] },
      {
        path: 'services',
        children: [
          { path: 'create', component: KtbCreateServiceComponent },
          { path: 'edit/:serviceName', component: KtbEditServiceComponent },
          { path: '', pathMatch: 'full', component: KtbServiceSettingsOverviewComponent },
        ],
      },
      {
        path: 'uniform',
        children: [
          {
            path: 'integrations',
            loadChildren: () =>
              import('./ktb-integration-view/ktb-integration-view.module').then((m) => m.KtbIntegrationViewModule),
          },
          {
            path: 'secrets',
            component: KtbSecretsListComponent,
          },
          { path: 'secrets/add', component: KtbCreateSecretFormComponent },
          { path: '', pathMatch: 'full', redirectTo: 'integrations' },
        ],
      },
      {
        path: 'support',
        children: [
          {
            path: 'common-use-cases',
            loadChildren: () =>
              import('./ktb-common-use-cases-view/ktb-common-use-cases-view.module').then(
                (m) => m.KtbCommonUseCasesViewModule
              ),
          },
          { path: '', pathMatch: 'full', redirectTo: 'common-use-cases' },
        ],
      },
      { path: '', pathMatch: 'full', redirectTo: 'project' },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbSettingsViewRoutingModule {}
