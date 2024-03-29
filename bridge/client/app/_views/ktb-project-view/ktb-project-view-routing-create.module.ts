import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PendingChangesGuard } from '../../_guards/pending-changes.guard';
import { NotFoundComponent } from '../../not-found/not-found.component';
import { KtbProjectViewComponent } from './ktb-project-view.component';
import { KtbProjectSettingsComponent } from '../ktb-settings-view/ktb-project-settings/ktb-project-settings.component';

const routes: Routes = [
  {
    path: '',
    component: KtbProjectViewComponent,
    data: { createMode: true },
    children: [
      {
        path: 'project',
        component: KtbProjectSettingsComponent,
        canDeactivate: [PendingChangesGuard],
      },
      { path: '**', component: NotFoundComponent },
    ],
  },
  { path: '**', component: NotFoundComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbProjectViewRoutingCreateModule {}
