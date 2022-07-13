import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { KtbProjectViewComponent } from './ktb-project-view.component';

const lazyLoadEnvironmentView = (): Promise<unknown> =>
  import('../ktb-environment-view/ktb-environment-view.module').then((m) => m.KtbEnvironmentViewModule);

const routes: Routes = [
  {
    path: '',
    component: KtbProjectViewComponent,
    children: [
      { path: '', pathMatch: 'full', loadChildren: lazyLoadEnvironmentView },
      { path: 'environment', loadChildren: lazyLoadEnvironmentView },
      { path: 'environment/stage/:stageName', loadChildren: lazyLoadEnvironmentView },
      {
        path: 'settings',
        loadChildren: () =>
          import('../ktb-settings-view/ktb-settings-view.module').then((m) => m.KtbSettingsViewModule),
      },
      {
        path: 'service',
        loadChildren: () => import('../ktb-service-view/ktb-service-view.module').then((m) => m.KtbServiceViewModule),
      },
      {
        path: 'sequence',
        loadChildren: () =>
          import('../ktb-sequence-view/ktb-sequence-view.module').then((m) => m.KtbSequenceViewModule),
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbProjectViewRoutingModule {}
